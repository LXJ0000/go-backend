package bootstrap

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"log"
	"os"
)

func NewAliyunClient(env *Env) *sms.Client {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	String := func(str string) *string {
		return &str
	}
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID")),
		//// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: String(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")),

		Endpoint: String(env.SMSEndpoint),
		RegionId: String(env.SMSRegionId),
	}
	_result, err := sms.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	return _result
}
