package models

import "gorm.io/gorm"

type InviteRelationships struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey" json:"id"`
	InviterID  uint   `gorm:"column:inviter_id" json:"inviter_id"`
	InviteeId  uint   `gorm:"column:invitee_id" json:"invitee_id"`
	InviteCode string `gorm:"column:invite_code" json:"invite_code"`
}
