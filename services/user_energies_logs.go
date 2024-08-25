package services

import (
	"ginDemo/database"
	"ginDemo/models"
)

func UpdateEnergiesLogs(userPointsLogs *models.UserEnergiesLogs) {
	database.DB.Create(userPointsLogs)
}
