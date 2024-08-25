package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username         string `gorm:"column:username" json:"username"`
	Points           string `gorm:"column:points; default:0" json:"points"`
	Energy           int    `gorm:"column:energy; default:5" json:"energy"`
	EarnEnergy       int    `gorm:"column:earn_energy" json:"earn_energy"`
	TelegramID       int    `gorm:"column:telegram_id" json:"telegram_id"`
	PhotoURL         string `gorm:"column:photo_url" json:"photo_url"`
	TelegramUsername string `gorm:"column:telegram_username" json:"telegram_username"`
	InviteCode       string `gorm:"column:invite_code" json:"invite_code"`
}
