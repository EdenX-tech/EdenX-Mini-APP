package services

import (
	"ginDemo/database"
	"ginDemo/models"
)

func SelectInviteeId(userId uint) (*models.InviteRelationships, error) {
	var inviteRelationships models.InviteRelationships
	if err := database.DB.Where("invitee_id = ?", userId).First(&inviteRelationships).Error; err != nil {
		return nil, err
	}
	return &inviteRelationships, nil
}

func UpdateInviteLogs(inviterId, inviteeId uint, inviteCode string) uint {
	var inviteRelationships = &models.InviteRelationships{
		InviterID:  inviterId,
		InviteeId:  inviteeId,
		InviteCode: inviteCode,
	}
	result := database.DB.Create(&inviteRelationships)

	if result.Error != nil {
		return 0
	} else {
		return inviteRelationships.ID
	}
}

func InviteList(inviterId uint, page, pageSize int) ([]*models.ResponseInviteDetail, error) {
	var inviteDetail []*models.ResponseInviteDetail
	println("inviterId:", inviterId)
	// 计算偏移量
	offset := (page - 1) * pageSize
	query := database.DB.Table("invite_relationships ir").
		Select("ir.invitee_id, u.username, u.photo_url, e.energies").
		Joins("join users u on ir.inviter_id = u.id").
		Joins("join user_energies_logs e on ir.id = e.data_id").
		Where("ir.inviter_id", inviterId).
		Order("ir.created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Scan(&inviteDetail)

	if query.Error != nil {
		return nil, query.Error
	}

	return inviteDetail, nil
}
