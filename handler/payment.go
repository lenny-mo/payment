package handler

import (
	"context"
	"errors"

	"github.com/lenny-mo/payment/domain/models"
	"github.com/lenny-mo/payment/domain/services"
	"github.com/lenny-mo/payment/proto/payment"

	"github.com/google/uuid"
)

// 实现下面pb的server api
//
//	type PaymentServiceHandler interface {
//		// 发起支付请求
//		// 获取订单ID，查阅order 表 获取 userid, user email,
//		// 根据order_item的ID和timestamp查询价格表 item_price, 相加每件商品价格得到总价
//		// 通过paypal查询用户是否有足够余额
//		// 如果充足则扣费并且交易记录要直接写入数据库，否则返回余额不足错误
//		MakePayment(context.Context, *MakePaymentRequest, *MakePaymentResponse) error
//		// 查询订单支付状态, 根据payment订单的ID
//		// 1. 接收GetPaymentStatusRequest中的支付订单ID。
//		// 2. 根据支付订单ID查询支付系统数据库，获取当前支付状态。
//		// 3. 考虑到支付可能是异步完成的，确保实时或定时查询支付渠道，获取最新状态。
//		// 4. 如果长时间未收到支付结果通知，触发主动查询流程，确保及时更新支付状态。
//		// 5. 返回支付状态，包括成功、失败、处理中等状态信息。
//		GetPaymentStatus(context.Context, *GetPaymentStatusRequest, *GetPaymentStatusResponse) error
//		// 更新支付信息
//		UpdatePayment(context.Context, *UpdatePaymentRequest, *UpdatePaymentResponse) error
//	}
type PaymentHandler struct {
	PaymentService services.PaymentService
}

func (p *PaymentHandler) MakePayment(ctx context.Context, req *payment.MakePaymentRequest, res *payment.MakePaymentResponse) error {
	// 构建一个payment 结构体
	paymentUUID := uuid.New().String()
	pdata := models.Payment{
		TransactionId:     paymentUUID,
		OrderId:           req.OrderId,
		TransactionStatus: 0, // 未支付
		PaymentMethod:     "paypal",
	}

	rowaffect, err := p.PaymentService.CreatePaymentRecord(pdata)
	if rowaffect == 0 || err != nil {
		res.Code = int32(FailedCode)
		res.CodeMsg = codeMsgMap[FailedCode]
		return err
	}

	res.Code = int32(SuccessCode)
	res.CodeMsg = codeMsgMap[SuccessCode]
	res.PaymentID = paymentUUID

	return nil
}

func (p *PaymentHandler) GetPaymentStatus(ctx context.Context, req *payment.GetPaymentStatusRequest, res *payment.GetPaymentStatusResponse) error {
	paymentdata, err := p.PaymentService.FindPaymentRecordById(req.PaymentId)
	if err != nil {
		return err
	}
	res.PaymentData = &payment.Payment{
		TransactionId:     paymentdata.TransactionId,
		OrderId:           paymentdata.OrderId,
		TransactionStatus: int32(paymentdata.TransactionStatus),
		PaymentMethod:     paymentdata.PaymentMethod,
	}
	return nil
}

func (p *PaymentHandler) UpdatePayment(ctx context.Context, req *payment.UpdatePaymentRequest, res *payment.UpdatePaymentResponse) error {
	// 根据req 的payment data 构建一个models.Payment
	reqPaymentData := models.Payment{
		TransactionId:     req.PaymentData.TransactionId,
		OrderId:           req.PaymentData.OrderId,
		PaymentMethod:     req.PaymentData.PaymentMethod,
		TransactionStatus: int8(req.PaymentData.TransactionStatus),
	}
	// 传递给service func
	rowAffected, err := p.PaymentService.UpdatePaymentRecord(reqPaymentData)
	if err != nil { // 更新出错
		return err
	}
	if rowAffected != 1 {
		return errors.New("update payment table more than 1 row")
	}
	res.PaymentId = req.PaymentData.TransactionId
	return nil
}
