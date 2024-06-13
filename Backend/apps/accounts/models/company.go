package models

import "gorm.io/gorm"
import base "GoTransact/apps/base"

type Company struct {
	gorm.Model
	base.Base   
	Name   string `gorm:"size:255"`
	UserID int    `gorm:"unique"`
	// User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
