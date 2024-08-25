package models

import "gorm.io/gorm"

type UserEnergiesLogs struct {
	gorm.Model
	UserID       uint `gorm:"column:user_id"`
	EnergiesType uint `gorm:"column:energies_type"`
	DataID       uint `gorm:"column:data_id"`
}

// 好友列表
type ResponseInviteDetail struct {
	InviteeID uint   `json:"invitee_id"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	Energies  int    `json:"energies"`
}
