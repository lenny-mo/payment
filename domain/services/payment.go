package services

// 业务代码

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/lenny-mo/payment/domain/dao"
	"github.com/lenny-mo/payment/domain/models"
	"github.com/lenny-mo/payment/middleware"
	"github.com/lenny-mo/payment/paypal"
)

type PaymentServiceInterface interface {
	// 返回payment 雪花算法ID
	CreatePaymentRecord(models.Payment) (int64, error)
	FindPaymentRecordById(string) (models.Payment, error)
	// 返回rowaffected
	UpdatePaymentRecord(models.Payment) (int64, error)
}

type chanData struct {
	ok  bool
	err error
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
// TODO: 添加熔断
func (p *PaymentService) CreatePaymentRecord(payment models.Payment) (int64, error) {

	// 1 处理支付渠道的异常故障或网络问题时，实施及时熔断
	// 获取paypal 的token
	accessToken := paypal.GenerateToken()
	// 生成一个payment request id
	paymentRequestId := paypal.UUID()
	// 计算订单的总金额
	amount := "10"
	// 再创建订单
	orderMapping := paypal.CreateOrder(accessToken, paymentRequestId, amount)
	fmt.Println(orderMapping)
	// 获取创建的订单id和需要用户支付的url
	orderId, userPayURL := orderMapping["id"], orderMapping["payer-action"]
	// 这里可以先尝试在控制台打印url
	fmt.Println("请点击付款：", userPayURL)
	// paypal 需要商家主动获取订单信息

	// 使用waitgroup + select
	var rowAffected int64

	ch := make(chan chanData, 1) // 创建一个通道 只能接收一个支付成功信号
	timer := time.NewTimer(15 * time.Minute)

	var wg sync.WaitGroup

outerloop:
	for {
		// 持续监听用户支付信息
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			ok, err := paypal.CapturePayment(orderId, paymentRequestId, accessToken)
			if ok {
				select {
				case ch <- chanData{ok: ok, err: err}:
					return
				default:
					return // ch阻塞，说明已经成功拿到支付信息, 退出
				}
			}
		}(&wg)

		select {
		case <-timer.C:
			// 等待15mins，如果倒计时到了，就返回支付失败, 结束所有的goroutine
			return 0, errors.New("user doesn't pay the order")
		case <-ch:
			// 如果成功，则把payment 先写入数据库，再更新redis缓存
			wg.Wait()                                                            // 等待所有goroutine都退出
			rowAffected, _ = p.Dao.CreatePaymentRecord(payment)                  // 存储进入数据库
			middleware.RedisSet(strconv.FormatInt(payment.OrderId, 10), payment) // 更新缓存
			break outerloop                                                      //跳出当前for循环
		default:
			fmt.Println("请点击付款：", userPayURL)
			time.Sleep(10 * time.Second)
			fmt.Println("waiting for user paying the order")
		}
	}

	// for {
	// 	if ok, _ := paypal.CapturePayment(orderId, paymentRequestId, accessToken); ok {
	// 		break
	// 	}
	// 	fmt.Println("请点击付款：", userPayURL)
	// 	time.Sleep(5 * time.Second)
	// }

	// var err error
	// rowAffected, err = p.Dao.CreatePaymentRecord(payment) // 存储进入数据库
	// if err != nil {
	// 	// zap
	// 	return 0, err
	// }
	// ok := middleware.RedisSet(strconv.FormatInt(payment.OrderId, 10), payment) // 更新缓存

	// if !ok {
	// 	fmt.Println("插入rediss失败")
	// }
	fmt.Println("receive payment from buyer")
	return rowAffected, nil
}

func (p *PaymentService) FindPaymentRecordById(paymentId string) (models.Payment, error) {
	// 1 先查询缓存
	var data string
	// 判断是否返回nil 值
	if v := middleware.RedisGet(paymentId); v != nil {
		data = v.(string)
	}
	if len(data) != 0 { // 缓存命中

		// 根据str反序列化成一个结构体
		var paymentdata models.Payment
		if err := json.Unmarshal([]byte(data), &paymentdata); err != nil {
			return models.Payment{}, err
		} else {
			return paymentdata, nil
		}
	}

	//2 缓存查询不到再查数据库，查到数据先写缓存再返回
	paymentdata, err := p.Dao.FindPaymentRecordById(paymentId)
	if err != nil {
		return models.Payment{}, err
	}

	middleware.RedisSet(paymentId, paymentdata)

	return paymentdata, nil
}

func (p *PaymentService) UpdatePaymentRecord(payment models.Payment) (int64, error) {
	// 延迟双删
	// 1 先删除缓存，更新数据库
	middleware.RedisSet(strconv.FormatInt(payment.OrderId, 10), payment)
	rowAffected, err := p.Dao.UpdatePaymentRecord(payment)
	if err != nil {
		return 0, err
	}
	// 2 再删除一次缓存, 容忍一定时间的脏数据
	middleware.RedisSet(strconv.FormatInt(payment.OrderId, 10), payment)
	return rowAffected, nil
}
