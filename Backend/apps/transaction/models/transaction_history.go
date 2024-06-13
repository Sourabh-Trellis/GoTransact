package models

import (
	base "GoTransact/apps/base"

	"gorm.io/gorm"
)

type TransactionHistory struct {
	gorm.Model
	base.Base
	TransactionID uint `gorm:"not null"`
	// Transaction   TransactionRequest `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status      TransactionStatus `gorm:"type:varchar(20);not null"`
	Description string            `gorm:"size:255"`
	Amount      float64           `gorm:"type:float"`
}
