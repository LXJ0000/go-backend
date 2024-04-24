package usecase

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-lib/queue"
)

type PostRankUsecase struct {
	interactionUsecase domain.InteractionUseCase
	postUsecase        domain.PostUsecase
	batchSize          int
	n                  int
	getScore           func(likeCnt int, updateTime time.Time) float64
}

func NewPostRankUsecase(interactionUsecase domain.InteractionUseCase, postUsecase domain.PostUsecase) *PostRankUsecase {
	return &PostRankUsecase{
		interactionUsecase: interactionUsecase,
		postUsecase:        postUsecase,
		batchSize:          100,
		n:                  100,
		getScore: func(likeCnt int, updateTime time.Time) float64 {
			var g float64 = 1.5
			cnt := float64(likeCnt)
			return (cnt - 1) / math.Pow(float64(updateTime.UnixMilli())+2, g)
		},
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
	}()
	return nil
}

func (ru *PostRankUsecase) topN(c context.Context) ([]domain.Post, error) {
	filter := map[string]interface{}{
		"status": domain.PostStatusPublish,
	}
	offset := 0
	type pair struct {
		post  domain.Post
		score float64
	}
	heap := queue.NewPriorityQueue[pair](func(first, second pair) bool {
		return first.score < second.score
	}) // 可以使用非并发安全 // TODO
	var minScore float64
	for {
		posts, err := ru.postUsecase.List(c, filter, offset, ru.batchSize)
		if err != nil {
			return nil, err
		}
		var postIDs []int64
		for _, post := range posts {
			postIDs = append(postIDs, post.PostID)
		}
		interactions, err := ru.interactionUsecase.GetByIDs(c, domain.BizPost, postIDs)
		if err != nil {
			return nil, err
		}
		// 合并计算 score
		for _, post := range posts {
			interaction := interactions[post.PostID]
			score := ru.getScore(interaction.LikeCnt, interaction.UpdatedAt)
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
		if len(posts) < ru.batchSize {
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
