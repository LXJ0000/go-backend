package ratelimit

import (
	"context"
	"github.com/LXJ0000/go-backend/internal/domain"
	sms2 "github.com/LXJ0000/go-backend/internal/usecase/sms"
	"github.com/LXJ0000/go-backend/utils/ratelimitutil"
)

type Service struct {
	service sms2.Service
	limiter ratelimitutil.Limiter
}

func NewService(service sms2.Service, limiter ratelimitutil.Limiter) sms2.Service {
	return &Service{
		service: service,
		limiter: limiter,
	}
}

func (s *Service) Send(ctx context.Context, templateID string, args []sms2.Param, numbers ...string) error {
	allow, err := s.limiter.Limit(ctx, "sms:limit")
	if err != nil {
		return err
	}
	if !allow {
		return domain.ErrRateLimit
	}
	return s.service.Send(ctx, templateID, args, numbers...)
}
