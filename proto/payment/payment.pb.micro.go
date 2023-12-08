// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/payment.proto

package payment

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/go-micro/v2/api"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for PaymentService service

func NewPaymentServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for PaymentService service

type PaymentService interface {
	// 发起支付请求
	// 获取订单ID，查阅order 表 获取 userid, user email,
	// 根据order_item的ID和timestamp查询价格表 item_price, 相加每件商品价格得到总价
	// 通过paypal查询用户是否有足够余额
	// 如果充足则扣费并且交易记录要直接写入数据库，否则返回余额不足错误
	MakePayment(ctx context.Context, in *MakePaymentRequest, opts ...client.CallOption) (*MakePaymentResponse, error)
	// 查询订单支付状态, 根据payment订单的ID
	// 1. 接收GetPaymentStatusRequest中的支付订单ID。
	// 2. 根据支付订单ID查询支付系统数据库，获取当前支付状态。
	// 3. 考虑到支付可能是异步完成的，确保实时或定时查询支付渠道，获取最新状态。
	// 4. 如果长时间未收到支付结果通知，触发主动查询流程，确保及时更新支付状态。
	// 5. 返回支付状态，包括成功、失败、处理中等状态信息。
	GetPaymentStatus(ctx context.Context, in *GetPaymentStatusRequest, opts ...client.CallOption) (*GetPaymentStatusResponse, error)
	// 更新支付信息
	UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, opts ...client.CallOption) (*UpdatePaymentResponse, error)
}

type paymentService struct {
	c    client.Client
	name string
}

func NewPaymentService(name string, c client.Client) PaymentService {
	return &paymentService{
		c:    c,
		name: name,
	}
}

func (c *paymentService) MakePayment(ctx context.Context, in *MakePaymentRequest, opts ...client.CallOption) (*MakePaymentResponse, error) {
	req := c.c.NewRequest(c.name, "PaymentService.MakePayment", in)
	out := new(MakePaymentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentService) GetPaymentStatus(ctx context.Context, in *GetPaymentStatusRequest, opts ...client.CallOption) (*GetPaymentStatusResponse, error) {
	req := c.c.NewRequest(c.name, "PaymentService.GetPaymentStatus", in)
	out := new(GetPaymentStatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentService) UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, opts ...client.CallOption) (*UpdatePaymentResponse, error) {
	req := c.c.NewRequest(c.name, "PaymentService.UpdatePayment", in)
	out := new(UpdatePaymentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PaymentService service

type PaymentServiceHandler interface {
	// 发起支付请求
	// 获取订单ID，查阅order 表 获取 userid, user email,
	// 根据order_item的ID和timestamp查询价格表 item_price, 相加每件商品价格得到总价
	// 通过paypal查询用户是否有足够余额
	// 如果充足则扣费并且交易记录要直接写入数据库，否则返回余额不足错误
	MakePayment(context.Context, *MakePaymentRequest, *MakePaymentResponse) error
	// 查询订单支付状态, 根据payment订单的ID
	// 1. 接收GetPaymentStatusRequest中的支付订单ID。
	// 2. 根据支付订单ID查询支付系统数据库，获取当前支付状态。
	// 3. 考虑到支付可能是异步完成的，确保实时或定时查询支付渠道，获取最新状态。
	// 4. 如果长时间未收到支付结果通知，触发主动查询流程，确保及时更新支付状态。
	// 5. 返回支付状态，包括成功、失败、处理中等状态信息。
	GetPaymentStatus(context.Context, *GetPaymentStatusRequest, *GetPaymentStatusResponse) error
	// 更新支付信息
	UpdatePayment(context.Context, *UpdatePaymentRequest, *UpdatePaymentResponse) error
}

func RegisterPaymentServiceHandler(s server.Server, hdlr PaymentServiceHandler, opts ...server.HandlerOption) error {
	type paymentService interface {
		MakePayment(ctx context.Context, in *MakePaymentRequest, out *MakePaymentResponse) error
		GetPaymentStatus(ctx context.Context, in *GetPaymentStatusRequest, out *GetPaymentStatusResponse) error
		UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, out *UpdatePaymentResponse) error
	}
	type PaymentService struct {
		paymentService
	}
	h := &paymentServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&PaymentService{h}, opts...))
}

type paymentServiceHandler struct {
	PaymentServiceHandler
}

func (h *paymentServiceHandler) MakePayment(ctx context.Context, in *MakePaymentRequest, out *MakePaymentResponse) error {
	return h.PaymentServiceHandler.MakePayment(ctx, in, out)
}

func (h *paymentServiceHandler) GetPaymentStatus(ctx context.Context, in *GetPaymentStatusRequest, out *GetPaymentStatusResponse) error {
	return h.PaymentServiceHandler.GetPaymentStatus(ctx, in, out)
}

func (h *paymentServiceHandler) UpdatePayment(ctx context.Context, in *UpdatePaymentRequest, out *UpdatePaymentResponse) error {
	return h.PaymentServiceHandler.UpdatePayment(ctx, in, out)
}
