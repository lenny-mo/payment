package models

import "gorm.io/gorm"

type Payment struct {
	// reference: proto/payment/.go文件 Payment struct
	TransactionId     int64  `gorm:"private_key;not_null"`
	OrderId           int64  `gorm:"not_null;unique"`
	TransactionStatus int8   `gorm:"not_null"`
	PaymentMethod     string `gorm:"not_null"`
	gorm.Model               // 自增id, createat
}

