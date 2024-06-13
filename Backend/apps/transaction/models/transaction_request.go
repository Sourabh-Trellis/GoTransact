package models

import (
	base "GoTransact/apps/base"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusProcessing TransactionStatus = "processing"
	StatusSuccess    TransactionStatus = "success"
	StatusFailed     TransactionStatus = "failed"
)

type TransactionRequest struct {
	gorm.Model
	base.Base
	UserID uint   `gorm:""`
	// User                   User              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status             TransactionStatus  `gorm:"type:varchar(20);not null;default:'pending'"`
	Payment_Gateway_id uint               `gorm:""`
	Description        string             `gorm:"size:255"`
	Amount             float64            `gorm:"type:float"`
	TransactionHistory TransactionHistory `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
