package usecase

import (
	"context"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-lib/queue"
	"log/slog"
	"math"
	"time"
)

type PostRankUsecase struct {
	interactionRepository domain.InteractionRepository
	postRepository        domain.PostRepository
	rankRepository        domain.RankRepository
	batchSize             int
	n                     int
	getScore              func(likeCnt int, updateTime time.Time) float64
	contextTimeout        time.Duration
}

func NewPostRankUsecase(
	interactionRepository domain.InteractionRepository,
	postRepository domain.PostRepository,
	rankRepository domain.RankRepository,
	contextTimeout time.Duration,
) *PostRankUsecase {
	return &PostRankUsecase{
		interactionRepository: interactionRepository,
		postRepository:        postRepository,
		rankRepository:        rankRepository,
		batchSize:             100,
		n:                     100,
		getScore: func(likeCnt int, updateTime time.Time) float64 {
			var g float64 = 1.5
			cnt := float64(likeCnt)
			t := time.Since(updateTime).Seconds()
			return (cnt - 1) / math.Pow(t+2, g)
		},
		contextTimeout: contextTimeout,
	}
}

func (ru *PostRankUsecase) TopN(c context.Context) error {
	posts, err := ru.topN(c)
	if err != nil {
		return err
	}
	slog.Info("topN", "posts", posts)
	go func() {
		// cache posts
		if err := ru.rankRepository.ReplaceTopN(context.Background(), posts, time.Minute*5); err != nil {
			slog.Error("Cache ReplaceTopN Failed", "Error", err.Error())
		} // TODO 配置expiration
	}()
	return nil
}

func (ru *PostRankUsecase) topN(c context.Context) ([]domain.Post, error) {
	now := time.Now()

	offset := 0
	type pair struct {
		post  domain.Post
		score float64
	}
	heap := queue.NewPriorityQueue(func(first, second pair) bool {
		return first.score < second.score
	}) // 可以使用非并发安全 // TODO
	var minScore float64
	for {
		posts, err := ru.postRepository.FindTopNPage(c, offset, ru.batchSize, now)
		if err != nil {
			return nil, err
		}
		var postIDs []int64
		for _, post := range posts {
			postIDs = append(postIDs, post.PostID)
		}
		interactions, err := ru.interactionRepository.GetByIDs(c, domain.BizPost, postIDs)
		if err != nil {
			return nil, err
		}
		// 合并计算 score
		for _, post := range posts {
			interaction := interactions[post.PostID]
			score := ru.getScore(interaction.LikeCnt, time.UnixMicro(interaction.UpdatedAt))
			//score := ru.getScore(interaction.LikeCnt, time.UnixMicro(interaction.UpdatedAt))

			// solve heap
			if heap.Size() < ru.n {
				heap.Push(pair{post: post, score: score})
				minScore = min(minScore, score)
				continue
			}
			if score > minScore { // 避免每次都 heap.Front()
				heap.Pop()
				heap.Push(pair{post: post, score: score})
				minScore = heap.Front().score
			}
		}
		// 不够一批 或者 时间跨度为7天，直接中断
		if len(posts) < ru.batchSize || now.Sub(time.UnixMicro(posts[len(posts)-1].UpdatedAt)).Hours() > 7*24 {
			break // 完啦
		}
		offset += ru.batchSize
	}
	res := make([]domain.Post, heap.Size())
	for i := heap.Size() - 1; heap.Size() > 0; i-- {
		post := heap.Pop().post
		res[i] = post
	}
	return res, nil
}

func (ru *PostRankUsecase) GetTopN(c context.Context) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()
	return ru.rankRepository.GetTopN(ctx)
}
