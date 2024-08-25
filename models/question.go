package models

type Question struct {
	ID           uint     `gorm:"primaryKey"`
	Text         string   `gorm:"type:varchar(255);not null"`
	QuestionType uint     `gorm:"not null"`
	Options      []Option `gorm:"foreignKey:QuestionID"`
}

type Option struct {
	ID          uint   `gorm:"primaryKey"`
	QuestionID  uint   `gorm:"not null"`
	OptionText  string `gorm:"type:varchar(255);not null"`
	IsCorrect   int    `gorm:"type:int;not null"`
	CorrectText string `gorm:"type:varchar(255)"`
}
