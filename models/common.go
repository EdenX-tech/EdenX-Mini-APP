package models

type LoginInput struct {
	ID         int    `json:"id"  binding:"required"`
	FirstName  string `json:"first_name"  binding:"required"`
	LastName   string `json:"last_name"  binding:"required"`
	Username   string `json:"username"  binding:"required"`
	PhotoURL   string `json:"photo_url"`
	AuthDate   string `json:"auth_date"  binding:"required"`
	Hash       string `json:"hash"  binding:"required"`
	InviteCode string `json:"invite_code"`
}

type ReturnLoginResult struct {
	Token      string
	Username   string
	PhotoURL   string
	Points     string
	Energy     int
	IsOpen     bool
	InviteCode string
}

type MessageTodoType struct {
	TodoType        string `json:"message_todo_type"`
	SelectedOptions uint   `json:"selected_option"`
	QuestionType    uint   `json:"question_type"`
}

type Answer struct {
	SelectedOptions uint `json:"selected_option"`
}

type QuestionResponse struct {
	ID      uint           `json:"id"`
	Text    string         `json:"text"`
	Options []QuestionData `json:"options"`
}

type CorrectOptionResponse struct {
	CorrectOption uint
	CorrectText   string
}

type AnswerResponse struct {
	CorrectOption uint   `json:"correct_option"`
	CorrectText   string `json:"correct_text"`
	Points        string `json:"points"`
	Energy        int    `json:"energy"`
	AwardsPoints  string `json:"awards_points"`
	IsCorrect     bool   `json:"is_correct"`
}

type UserQuestion struct {
	UserID        uint
	QuestionID    uint
	CorrectOption uint
	CorrectText   string
}

type QuestionData struct {
	OptionID   uint
	OptionText string
}

type ErrorResponse struct {
	ErrorCode uint   `json:"error_code"`
	Message   string `json:"message"`
}

type Response struct {
	ErrorCode uint        `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

type InviteInput struct {
	Page    int `json:"page"  binding:"required"`
	PerPage int `json:"per_page"  binding:"required"`
}

type EncryptedInput struct {
	Data string `json:"data"`
}

type AuthTelegramInput struct {
	InitData string `json:"initData"`
}

type WsAuthInput struct {
	Authorization  string `json:"authorization"`
	AcceptLanguage string `json:"accept_language"`
}

type PageAndPerPage struct {
	Page    int `json:"page" binding:"required"`
	PerPage int `json:"per_page" binding:"required"`
}
