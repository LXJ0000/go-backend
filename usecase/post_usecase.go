package usecase

import (
	"github.com/LXJ0000/go-backend/domain"
	"github.com/LXJ0000/go-backend/event"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

type postUsecase struct {
	repo           domain.PostRepository
	contextTimeout time.Duration
	producer       event.Producer
}

func NewPostUsecase(repo domain.PostRepository, timeout time.Duration, producer event.Producer) domain.PostUsecase {
	return &postUsecase{
		repo:           repo,
		contextTimeout: timeout,
		producer:       producer,
	}
}
func (uc *postUsecase) Create(c context.Context, post *domain.Post) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.Create(ctx, post)
}

func (uc *postUsecase) List(c context.Context, filter interface{}, page, size int) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	return uc.repo.FindPage(ctx, filter, page, size)
}

func (uc *postUsecase) Info(c context.Context, postID int64) (domain.Post, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()
	post, err := uc.repo.GetByID(ctx, postID)
	if err == nil {
		go func() {
			// TODO context
			if err := uc.producer.ProduceReadEvent(c, event.ReadEvent{
				PostID: post.PostID,
				UserID: post.AuthorID,
			}); err != nil {
				slog.Warn("ProduceReadEvent Fail", "err", err.Error(), "UserID", post.AuthorID, "PostID", post.PostID)
			}
		}()
	}
	return post, err
}

//func (uc *postUsecase) ReplaceTopN(c context.Context, items []domain.Post, expiration time.Duration) error {
//	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
//	defer cancel()
//	return uc.repo.ReplaceTopN(ctx, items, expiration)
//}
