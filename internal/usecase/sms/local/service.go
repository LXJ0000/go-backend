package local

import (
	"context"
	sms2 "github.com/LXJ0000/go-backend/internal/usecase/sms"
)

type Service struct {
	appID    string
	signName string
}

func NewService() sms2.Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, templateID string, args []sms2.Param, numbers ...string) error {
	return nil
}
