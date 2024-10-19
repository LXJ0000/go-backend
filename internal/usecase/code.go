package usecase

import (
	"context"
	"fmt"
	"github.com/LXJ0000/go-backend/internal/domain"
	sms2 "github.com/LXJ0000/go-backend/internal/usecase/sms"
	"log/slog"
	"math/rand"
)

type CodeService struct {
	codeRepo domain.CodeRepository
	sms      sms2.Service
}

func NewCodeService(codeRepo domain.CodeRepository, sms sms2.Service) *CodeService {
	return &CodeService{codeRepo: codeRepo, sms: sms}
}

func (s *CodeService) Send(ctx context.Context, biz, number string) error {
	code := s.genCode()
	if err := s.codeRepo.SetCode(ctx, biz, number, code); err != nil {
		slog.Error("set code error", "error", err.Error(), "biz", biz, "number", number, "code", code)
		return err
	}
	if err := s.sms.Send(ctx, "SMS_474870192", []sms2.Param{{Name: "code", Value: code}}, number); err != nil {
		// redis set 成功，sms 发送失败 不能刪除 redis key 因为错误有可能是超时错误... 即短信发送成功，但是返回超时
		// 解决方案一：重试 让调用者自己决定重试方案 即sms 缺陷：用户重复收到验证码；一直重复一直失败，系统负载高
		// 解决方案二：
		slog.Error("send sms error", "error", err.Error(), "biz", biz, "number", number, "code", code)
		return err
	}
	return nil
}

func (s *CodeService) Verify(ctx context.Context, biz, number, code string) (bool, error) {
	return s.codeRepo.VerifyCode(ctx, biz, number, code)
}

func (s *CodeService) genCode() string {
	// 生成6位數隨機驗證碼 0 - 999999
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
