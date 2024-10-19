package tencent

import (
	"context"
	"fmt"
	sms2 "github.com/LXJ0000/go-backend/internal/usecase/sms"

	"github.com/LXJ0000/go-lib/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appID    string
	signName string
	client   *sms.Client
}

func NewService(appID, signName string, client *sms.Client) *Service {
	return &Service{
		appID:    appID,
		signName: signName,
		client:   client,
	}
}

func (s *Service) Send(ctx context.Context, templateID string, args []sms2.Param, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = &s.appID
	req.SignName = &s.signName
	req.TemplateId = &templateID
	req.PhoneNumberSet = slice.Map(numbers, func(number string) *string {
		return &number
	})
	req.TemplateParamSet = slice.Map(args, func(arg sms2.Param) *string {
		return &arg.Value
	})
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code != nil || *(status.Code) != "Ok" {
			return fmt.Errorf("send sms failed: %s, %s", *(status.Code), *(status.Message))
		}
	}
	return nil
}
