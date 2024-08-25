package services

import (
	"ginDemo/database"
	"ginDemo/models"
)

func UpdatePointsLogs(userEnergiesLogs *models.UserPointsLogs) {
	database.DB.Save(userEnergiesLogs)
}
