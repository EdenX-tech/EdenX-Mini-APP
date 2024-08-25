package services

import (
	"ginDemo/database"
	"ginDemo/models"
)

func GetRandomQuestion(question_type uint) (*models.Question, error) {
	var question models.Question
	err := database.DB.Where("question_type = ?", question_type).Order("RAND()").First(&question).Error
	if err != nil {
		return nil, err
	}

	err = database.DB.Where("question_id = ?", question.ID).Find(&question.Options).Error

	if err != nil {
		return nil, err
	}

	return &question, nil
}
