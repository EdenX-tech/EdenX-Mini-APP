package websocket

import (
	"encoding/json"
	"errors"
	"ginDemo/common"
	"ginDemo/models"
	"ginDemo/services"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var userQuestionCache = struct {
	sync.RWMutex
	m map[uint]models.UserQuestion
}{m: make(map[uint]models.UserQuestion)}

var activeSessions sync.Map

func sendQuestion(userID, questionType uint) (error, *models.QuestionResponse) {
	question, err := services.GetRandomQuestion(questionType)

	if err != nil {
		return err, nil
	}

	var options []models.QuestionData
	var correctOption uint
	var correctText string
	for _, option := range question.Options {
		options = append(options, models.QuestionData{OptionID: option.ID, OptionText: option.OptionText})
		if option.IsCorrect == 1 {
			correctOption = option.ID
			correctText = option.CorrectText
		}
	}
	CacheUserQuestion(userID, question.ID, correctOption, correctText)

	questionResponse := models.QuestionResponse{
		ID:      question.ID,
		Text:    question.Text,
		Options: options,
	}

	return nil, &questionResponse
}

func CacheUserQuestion(userID uint, questionID uint, correctOption uint, correctText string) {
	userQuestionCache.Lock()
	userQuestionCache.m[userID] = models.UserQuestion{
		UserID:        userID,
		QuestionID:    questionID,
		CorrectOption: correctOption,
		CorrectText:   correctText,
	}
	userQuestionCache.Unlock()
}

func WebSocketHandler(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		common.ErrorJson(1000, c)
		return
	}
	currentUser := user.(*models.User)

	if conn, ok := activeSessions.Load(currentUser.ID); ok {
		common.ErrorSocketJson(conn.(*websocket.Conn), 1003, c)
		conn.(*websocket.Conn).Close()
		activeSessions.Delete(currentUser.ID)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.ErrorJson(1000, c)
		return
	}

	done := make(chan struct{})

	activeSessions.Store(currentUser.ID, conn)

	defer func() {
		activeSessions.Delete(currentUser.ID)
		if err := conn.Close(); err != nil {
		}
		close(done)
	}()

	if err = common.SocketJson(conn, "Connection Successful", 0); err != nil {
		common.HandleError(conn, err, 1000, c)
		return
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
				time.Sleep(10 * time.Second)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			_, message, err := conn.ReadMessage()
			if common.HandleError(conn, err, 1001, c) {

				return
			}

			var encryptedInput models.EncryptedInput
			err = json.Unmarshal(message, &encryptedInput)
			if common.HandleError(conn, err, 1000, c) {

				return
			}

			decryptedJson, err := utils.Decrypt(encryptedInput.Data)
			if err != nil {
				if common.HandleError(conn, err, 1001, c) {
					return
				}
			}
			var messageTodoType models.MessageTodoType

			if err := json.Unmarshal(decryptedJson, &messageTodoType); err != nil {
				common.ErrorJson(1000, c)
				return
			}
			userEnergy, _ := services.GetUserByID(currentUser.ID)

			switch messageTodoType.TodoType {
			case viper.GetString("executionCondition.start"):
				energyAmount := userEnergy.Energy + userEnergy.EarnEnergy
				if energyAmount <= 0 {
					common.ErrorSocketJson(conn, 6000, c)
					return
				}

				errBool, questionResponse := TodoStartAnswering(currentUser.ID, messageTodoType.QuestionType)
				if errBool != true {
					common.ErrorSocketJson(conn, 6100, c)
					return
				}

				if err = common.SocketJson(conn, questionResponse, 0); err != nil {
					common.HandleError(conn, err, 1000, c)
					return
				}
				break
			case viper.GetString("executionCondition.end"):
				if err, response, errorCode := todoEndAnswering(userEnergy, messageTodoType.SelectedOptions); err != nil {
					common.HandleError(conn, err, errorCode, c)
					return
				} else {
					if err = common.SocketJson(conn, response, errorCode); err != nil {
						common.HandleError(conn, err, errorCode, c)
						return
					}
				}
				break
			default:
				common.ErrorSocketJson(conn, 1002, c)
				return
			}
		}
	}
}

func TodoStartAnswering(userID, questionType uint) (bool, *models.QuestionResponse) {
	err, questionResponse := sendQuestion(userID, questionType)
	if err != nil {
		return false, nil
	}

	return true, questionResponse
}

func todoEndAnswering(userEnergy *models.User, selectedOptions uint) (error, *models.AnswerResponse, uint) {
	configPoints := viper.GetString("awards.points")
	configPointsInt, _ := common.StrToInt(configPoints)

	userQuestionCache.RLock()
	userQuestion, exist := userQuestionCache.m[userEnergy.ID]
	userQuestionCache.RUnlock()
	errorCode := uint(0)

	if !exist {
		return errors.New("user question not found"), nil, 6100
	}

	response := models.AnswerResponse{
		CorrectOption: userQuestion.CorrectOption,
		CorrectText:   userQuestion.CorrectText,
	}

	if selectedOptions == userQuestion.CorrectOption {
		pointsInt, _ := strconv.Atoi(userEnergy.Points)
		pointsInt = pointsInt + configPointsInt
		userEnergy.Points = strconv.Itoa(pointsInt)
		response.AwardsPoints = configPoints
		response.IsCorrect = true
		services.UpdatePointsLogs(&models.UserPointsLogs{
			UserID: userEnergy.ID,
		})
	} else {
		if userEnergy.EarnEnergy >= 1 {
			userEnergy.EarnEnergy -= 1
		} else {
			userEnergy.Energy -= 1
		}
		response.AwardsPoints = "0"
		response.IsCorrect = false
		services.UpdateEnergiesLogs(&models.UserEnergiesLogs{
			UserID:       userEnergy.ID,
			EnergiesType: 1,
		})
	}
	err := services.UpdateUser(userEnergy)
	if err != nil {
		return err, nil, 6100
	}

	response.Points = userEnergy.Points
	response.Energy = userEnergy.Energy + userEnergy.EarnEnergy

	return nil, &response, errorCode
}
