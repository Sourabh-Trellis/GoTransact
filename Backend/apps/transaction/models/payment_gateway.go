package models

import (
	base "GoTransact/apps/base"

	"gorm.io/gorm"
)

type Payment_Gateway struct {
	gorm.Model
	base.Base
	Slug               string             `gorm:"size:255;unique"`
	Label              string             `gorm:"size:255"`
	TransactionRequest TransactionRequest `gorm:"foreignKey:Payment_Gateway_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
