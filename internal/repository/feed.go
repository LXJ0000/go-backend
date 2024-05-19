package repository

import "github.com/LXJ0000/go-backend/internal/domain"

type feedRepository struct {
}

func NewFeedRepository() domain.FeedRepository {
	return &feedRepository{}
}
