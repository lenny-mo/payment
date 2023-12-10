package dao

import (
	"github.com/lenny-mo/payment/domain/models"

	"gorm.io/gorm"
)

type PaymentDAOInterface interface {
	// 返回payment 雪花算法ID
	CreatePaymentRecord(models.Payment) (int64, error)
	FindPaymentRecordById(string) (models.Payment, error)
	// 返回rowaffected
	UpdatePaymentRecord(models.Payment) (int64, error)
}

type PaymentDAO struct {
	db *gorm.DB
}

func NewPaymentDAO(db *gorm.DB) PaymentDAOInterface {
	return &PaymentDAO{
		db: db,
	}
}

func (p *PaymentDAO) CreatePaymentRecord(data models.Payment) (int64, error) {
	result := p.db.Create(data)
	return result.RowsAffected, result.Error
}

func (p *PaymentDAO) FindPaymentRecordById(PaymentId string) (models.Payment, error) {
	data := new(models.Payment)
	result := p.db.First(data, "transaction_id = ?", PaymentId)
	if result.Error != nil {
		return models.Payment{}, result.Error
	}
	return *data, nil
}

func (p *PaymentDAO) UpdatePaymentRecord(payment models.Payment) (int64, error) {
	data, err := p.FindPaymentRecordById(payment.TransactionId)
	if err != nil {
		// 插入数据
		result := p.db.Create(payment)
		return result.RowsAffected, result.Error
	}

	// 找到数据之后，更新数据
	data.PaymentMethod = payment.PaymentMethod
	data.TransactionStatus = payment.TransactionStatus
	result := p.db.Save(data)
	return result.RowsAffected, result.Error
}
