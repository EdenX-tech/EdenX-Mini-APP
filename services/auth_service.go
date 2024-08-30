package services

import (
	"ginDemo/database"
	"ginDemo/models"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User

	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(userId uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	return database.DB.Save(user).Error
}

func GetUserTelegramId(telegramID int) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserInviteCode(inviteCode string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("invite_code = ?", inviteCode).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func AddEnergy(userID uint, increment int) {
	var user models.User

	database.DB.Model(&user).Where("id = ?", userID).UpdateColumn("earn_energy", gorm.Expr("earn_energy + ?", increment))

}

func RestoreEnergy() error {
	var user models.User
	if err := database.DB.Model(&user).Where("energy < ?", 5).UpdateColumn("energy", 5).Error; err != nil {
		println("energy error:", err.Error())
		return err
	}
	return nil
}
