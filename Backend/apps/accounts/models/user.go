package models

import (
	base "GoTransact/apps/base"
	"GoTransact/apps/transaction/models"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	base.Base
	FirstName          string                    `gorm:"size:255"`
	LastName           string                    `gorm:"size:255"`
	Email              string                    `gorm:"size:255;unique" `
	Password           string                    `gorm:"size:255"`
	Company            Company                   `gorm:"foreignKey:UserID"`
	TransactionRequest models.TransactionRequest `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
