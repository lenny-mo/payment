package services
// 业务代码

import (
	"payment/domain/dao"
	"payment/domain/models"
)

type PaymentServiceInterface interface {
	// 返回payment 雪花算法ID
	CreatePaymentRecord(models.Payment) (int64, error)
	FindPaymentRecordById(int64) (models.Payment, error)
	// 返回rowaffected
	UpdatePaymentRecord(models.Payment) (int64, error)
}

type PaymentService struct {
	Dao dao.PaymentDAO
}

func NewPaymentService(dao dao.PaymentDAO) PaymentServiceInterface {
	return &PaymentService{Dao: dao}
}

// CreatePaymentRecord 创建支付记录的方法。
//
// 此函数用于PaymentService服务层中，用于处理创建支付订单的业务逻辑。它接收一个Payment类型的对象，
// 并返回一个int64类型的支付记录ID和一个可能出现的error。
//
// 参数:
//
//	payment models.Payment - 一个Payment对象，包含了支付订单所需的所有信息，如用户ID，金额，支付方式等。
//
// 返回值:
//
//	int64 - 成功创建的支付记录的ID。这个ID是支付订单的唯一标识，可以用于后续的查询、更新等操作。
//	error - 如果在创建支付记录的过程中发生错误，则返回相应的错误信息。如果没有错误发生，则返回nil。
//
// 流程说明:
//  1. 验证payment对象中的数据是否完整和有效。
//  2. 将payment对象的数据保存到数据库中。
//  3. 如果数据库操作成功，返回新创建的支付记录的ID。
//  4. 如果在任何步骤中发生错误，捕获错误并返回。
//
// 注意:
//   - 此函数不直接处理与支付网关的交互，它只负责处理与支付记录相关的内部逻辑。
//   - 需要确保传入的payment对象符合业务规则和数据完整性要求。
//   - 函数实现应考虑到事务性，确保数据的一致性和完整性。
//   - 考虑到支付订单的超时问题，可以在此函数中实现MQ延时队列逻辑，当订单超时未支付时自动关闭订单。
//     这部分逻辑可以通过发送一个延时消息到MQ，该消息包含订单ID和超时时间。当消息被消费时，
//     检查订单状态并根据需要更新订单状态为“已关闭”。
//
// 示例:
//
//	paymentRecord := models.Payment{UserID: "1234", Amount: 100.00, Method: "CreditCard"}
//	recordID, err := paymentService.CreatePaymentRecord(paymentRecord)
//	if err != nil {
//	    // 处理错误
//	}
//	// 使用recordID进行后续操作
func (p *PaymentService) CreatePaymentRecord(payment models.Payment) (int64, error) {

	// 1 处理支付渠道的异常故障或网络问题时，实施及时熔断

	// 2 支付订单超时关闭 采用MQ延时队列（如RabbitMQ）处理超时订单

	// 3 支付结果通知上游使用延时重试队列

	return 0, nil
}

func (p *PaymentService) FindPaymentRecordById(paymentId int64) (models.Payment, error) {
	// 1 先查询缓存

	// 2 缓存查询不到再查数据库，查到数据先写缓存再返回

	return models.Payment{}, nil
}

func (p *PaymentService) UpdatePaymentRecord(payment models.Payment) (int64, error) {
	// 延迟双删除
	// 1 先删除缓存，更新数据库

	// 2 再删除一次缓存

	return 0, nil
}
