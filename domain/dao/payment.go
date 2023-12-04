package dao

import (
	"payment/domain/models"

	"gorm.io/gorm"
)

type PaymentDAOInterface interface {
	// 返回payment 雪花算法ID
	CreatePaymentRecord(models.Payment) (int64, error)
	FindPaymentRecordById(int64) (models.Payment, error)
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
	return 0, nil
}

func (p *PaymentDAO) FindPaymentRecordById(PaymentId int64) (models.Payment, error) {
	return models.Payment{}, nil
}

func (p *PaymentDAO) UpdatePaymentRecord(payment models.Payment) (int64, error) {
	return 0, nil
}
