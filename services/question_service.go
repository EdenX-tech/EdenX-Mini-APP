package services

import (
	"ginDemo/database"
	"ginDemo/models"
)

func GetRandomQuestion(questionType uint, language string) (*models.Question, error) {
	var question models.Question

	var questionField string
	var optionField string
	if language == "CN" {
		questionField = "id, cn_text AS text"
		optionField = "id, cn_option_text AS option_text, cn_correct_text as correct_text, is_correct, question_id AS question_id"
	} else {
		questionField = "id, text AS text"
		optionField = "id, option_text AS option_text, correct_text as correct_text, is_correct, question_id AS question_id"
	}
	err := database.DB.Select(questionField).Where("question_type = ?", questionType).Order("RAND()").First(&question).Error
	if err != nil {
		return nil, err
	}

	err = database.DB.Select(optionField).Where("question_id = ?", question.ID).Find(&question.Options).Error

	if err != nil {
		return nil, err
	}

	return &question, nil
}
