package models

import "gorm.io/gorm"

type UserPointsLogs struct {
	gorm.Model
	UserID uint `gorm:"column:user_id"`
}
