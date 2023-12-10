package utils

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
	"go.uber.org/zap"
)

// 客户端熔断模块
type clientWrapper struct {
	client.Client
}

// 这个函数名为 `Call`，是一个实现了 `client.Wrapper` 接口的方法，用于创建一个熔断器包装的客户端。
func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// 使用 hystrix 包中的 Do 方法，对请求进行熔断处理。
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		// 在正常执行时，打印请求的服务和端点信息。
		zap.L().Info(req.Service() + "." + req.Endpoint())
		// 调用原始客户端的 Call 方法执行请求。
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		// 处理熔断时的错误情况，并打印错误信息。
		zap.L().Error(e.Error())
		return e
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
