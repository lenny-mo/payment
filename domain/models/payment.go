package models

import "gorm.io/gorm"

type Payment struct {
	// reference: proto/payment/.go文件 Payment struct
	TransactionId     string `gorm:"not_null;unique;column:transaction_id;comment:'uuid'"`
	OrderId           int64  `gorm:"not_null;unique;column:order_id;comment:'order id'"` // 需要建立外键约束
	TransactionStatus int8   `gorm:"not_null;column:transaction_status;comment:'是否支付'"`
	PaymentMethod     string `gorm:"not_null;column:payment_method"`
	gorm.Model               // 自增id, createat
}
