package search

import (
	"context"
	"github.com/LXJ0000/go-backend/internal/domain"
	"golang.org/x/sync/errgroup"
	"strings"
)

type searchUsecase struct {
	repo domain.SearchRepository
}

func NewSearchUsecase() domain.SearchUsecase {
	return &searchUsecase{}
}

func (s *searchUsecase) Search(ctx context.Context, userID int64, cmd string) error {
	keywords := strings.Split(cmd, " ")
	// TODO 高级表达式cmd 过滤敏感词汇
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		users, err := s.repo.SearchUser(ctx, keywords...)
		res.Users = users
		return err
	})
	eg.Go(func() error {
		post, err := s.repo.SearchPost(ctx, userID, keywords...)
		res.Posts = post
		return err
	})
	// todo 添加不同类目 参考B站
	return eg.Wait()
}
