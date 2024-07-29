package repository

import (
	"context"
	"encoding/json"
	"github.com/LXJ0000/go-backend/internal/domain"
	"github.com/olivere/elastic/v7"
	"strings"
)

type searchRepository struct {
	client *elastic.Client
}

func NewSearchRepository() domain.SearchRepository {
	return &searchRepository{}
}

func (s *searchRepository) SearchUser(ctx context.Context, keywords ...string) ([]domain.User, error) {
	query := strings.Join(keywords, " ")
	// match nickname
	resp, err := s.client.Search(domain.ESUserIndex).
		Query(elastic.NewMatchQuery("nick_name", query)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]domain.User, 0, resp.Hits.TotalHits.Value)
	for _, hit := range resp.Hits.Hits {
		var user domain.User
		if err := json.Unmarshal(hit.Source, &user); err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return res, nil
}
func (s *searchRepository) SearchPost(ctx context.Context, userID int64, keywords ...string) ([]domain.Post, error) {
	query := strings.Join(keywords, " ")
	// match title or content
	// status 精确查找
	status := elastic.NewTermsQuery("status", domain.PostStatusPublish)
	// title
	title := elastic.NewMatchQuery("title", query)
	content := elastic.NewMatchQuery("content", query)
	titleOrContent := elastic.NewBoolQuery().Should(title, content)
	titleOrContentAndStatus := elastic.NewBoolQuery().Must(titleOrContent, status)

	resp, err := s.client.Search(domain.ESPostIndex).Query(titleOrContentAndStatus).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Post, 0, resp.Hits.TotalHits.Value)
	for _, hit := range resp.Hits.Hits {
		var post domain.Post
		if err := json.Unmarshal(hit.Source, &post); err != nil {
			return nil, err
		}
		res = append(res, post)
	}
	return res, nil
}
