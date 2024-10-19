package usecase

import (
	"log/slog"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/LXJ0000/go-backend/internal/event"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

type postUsecase struct {
	repo           domain.PostRepository
	contextTimeout time.Duration
	producer       event.Producer

	postRankUsecase *PostRankUsecase
}

func NewPostUsecase(repo domain.PostRepository, timeout time.Duration, producer event.Producer, postRankUsecase *PostRankUsecase) domain.PostUsecase {
	return &postUsecase{
		repo:            repo,
		contextTimeout:  timeout,
		producer:        producer,
		postRankUsecase: postRankUsecase,
	}
}

func (uc *postUsecase) Create(c context.Context, post *domain.Post) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Create(ctx, post)
}

func (uc *postUsecase) List(c context.Context, filter interface{}, page, size int) ([]domain.Post, int64, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	var (
		items []domain.Post
		cnt   int64
	)
	eg := errgroup.Group{}
	eg.Go(func() error {
		var err error
		items, err = uc.repo.List(ctx, filter, page, size)
		return err
	})
	eg.Go(func() error {
		var err error
		cnt, err = uc.repo.Count(ctx, filter)
		return err
	})
	if eg.Wait() != nil {
		return nil, 0, eg.Wait()
	}
	return items, cnt, nil
}

func (uc *postUsecase) Info(c context.Context, postID int64) (domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	post, err := uc.repo.GetByID(ctx, postID)
	if err == nil {
		go func() {
			// TODO context
			if err := uc.producer.ProduceReadEvent(context.Background(), event.ReadEvent{
				PostID: post.PostID,
				UserID: post.AuthorID,
			}); err != nil {
				slog.Warn("ProduceReadEvent Fail", "err", err.Error(), "UserID", post.AuthorID, "PostID", post.PostID)
			}
		}()
	}
	return post, err
}

func (uc *postUsecase) Count(c context.Context, filter interface{}) (int64, error) {
	return uc.repo.Count(c, filter)

}

// func (uc *postUsecase) ReplaceTopN(c context.Context, items []domain.Post, expiration time.Duration) error {
// 	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
// 	defer cancel()
// 	return uc.repo.ReplaceTopN(ctx, items, expiration)
// }

func (uc *postUsecase) TopN(c context.Context) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.postRankUsecase.GetTopN(ctx)
}
