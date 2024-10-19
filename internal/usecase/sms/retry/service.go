package retry

//import (
//	"context"
//	"github.com/LXJ0000/go-backend/internal/usecase/sms"
//)
//
//type Service struct {
//	svc sms.Service
//	cnt int
//}
//
//func (s *Service) Send(ctx context.Context, templateID string, args []sms.Param, numbers ...string) error {
//	var err error
//	err = s.svc.Send(ctx, templateID, args, numbers...)
//	if err != nil && s.cnt < 3 {
//		s.cnt++
//		err = s.svc.Send(ctx, templateID, args, numbers...)
//	}
//	return nil
//}
