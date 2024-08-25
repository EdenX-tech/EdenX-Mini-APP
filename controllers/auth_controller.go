package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"ginDemo/common"
	"ginDemo/models"
	"ginDemo/services"
	"ginDemo/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/url"
	"strings"
)

func Register(input *models.LoginInput) (*models.User, error) {
	username := input.FirstName + input.LastName
	inviteCode := common.GenerateInviteCode(8)

	user := models.User{
		Username:         username,
		TelegramID:       input.ID,
		PhotoURL:         input.PhotoURL,
		TelegramUsername: input.Username,
		InviteCode:       inviteCode,
	}

	if err := services.CreateUser(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func Login(c *gin.Context) {
	var encryptedInput models.EncryptedInput
	if err := c.ShouldBindJSON(&encryptedInput); err != nil {
		common.ErrorJson(1000, c)
		return
	}

	decryptedJson, err := utils.Decrypt(encryptedInput.Data)
	if err != nil {
		common.ErrorJson(1000, c)
		return
	}

	var authTelegramInput models.AuthTelegramInput

	if err := json.Unmarshal(decryptedJson, &authTelegramInput); err != nil {
		common.ErrorJson(1000, c)
		return
	}

	dataString, isValidata := validataLoginData(&authTelegramInput)
	if !isValidata {
		common.ErrorJson(1000, c)
		return
	}

	input, err := formatUserInfo(dataString)
	if err != nil {
		common.ErrorJson(1000, c)
		return
	}

	user, err := services.GetUserTelegramId(input.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user, err = Register(input)
			if err != nil {
				common.ErrorJson(1000, c)
				return
			}
			if input.InviteCode != "" {
				if checkBool := checkInviteeId(user.ID); checkBool {
					inviterId, err := services.GetUserInviteCode(input.InviteCode)
					if err != nil {
						common.ErrorJson(3001, c)
						return
					}
					dataId := services.UpdateInviteLogs(inviterId.ID, user.ID, input.InviteCode)
					services.AddEnergy(inviterId.ID, 1)
					services.UpdateEnergiesLogs(&models.UserEnergiesLogs{
						UserID:       inviterId.ID,
						EnergiesType: 3,
						DataID:       dataId,
					})

				}

			}
		}
	}

	token, tokenError := utils.GenerateToken(user.ID)
	if tokenError != nil {
		common.ErrorJson(1000, c)
		return
	}
	userEnergy := user.Energy + user.EarnEnergy

	returnLoginResult := models.ReturnLoginResult{
		Token:      token,
		Username:   user.Username,
		PhotoURL:   user.PhotoURL,
		Points:     user.Points,
		Energy:     userEnergy,
		IsOpen:     user.Energy > 0,
		InviteCode: user.InviteCode,
	}

	common.WbeJson(returnLoginResult, c)
}

func validataLoginData(initData *models.AuthTelegramInput) (string, bool) {
	var token = viper.GetString("telegram.botToken")
	pairs := strings.Split(initData.InitData, "&")

	var hash string
	var data []string

	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		key := kv[0]
		value, _ := url.QueryUnescape(kv[1]) // 对值进行解码
		if key == "hash" {
			hash = kv[1]
		} else {
			data = append(data, key+"="+value)
		}
	}
	joinedPairs := strings.Join(data, "\n")
	secretKey := generateSecretKey(token)
	calculatedHash := calculateHMAC(secretKey, joinedPairs)
	if calculatedHash == hash {
		return joinedPairs, true
	}
	return "", false
}

func formatUserInfo(data string) (*models.LoginInput, error) {
	start := strings.Index(data, "user=")

	userPart := data[start+len("user="):]

	end := strings.Index(userPart, "}")

	userJsonStr := userPart[:end+1]

	startParamStart := strings.Index(data, "start_param=")
	if startParamStart != -1 {
		startParamPart := data[startParamStart+len("start_param="):]

		startParamEnd := strings.IndexAny(startParamPart, "& \n")
		if startParamEnd == -1 {
			startParamEnd = len(startParamPart)
		}

		startParam := startParamPart[:startParamEnd]

		startParam = strings.TrimSpace(startParam)
		startParam = strings.ReplaceAll(startParam, `"`, `\"`)

		if startParam != "" && !strings.Contains(userJsonStr, `"start_param"`) {
			userJsonStr = userJsonStr[:len(userJsonStr)-1] + `,"start_param":"` + startParam + `"}`
		}

	}

	var input models.LoginInput

	if err := json.Unmarshal([]byte(userJsonStr), &input); err != nil {
		return nil, err
	}

	return &input, nil
}

func generateSecretKey(botToken string) []byte {
	mac := hmac.New(sha256.New, []byte("WebAppData"))
	mac.Write([]byte(botToken))
	return mac.Sum(nil)
}

func calculateHMAC(secretKey []byte, joinedPairs string) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(joinedPairs))
	return hex.EncodeToString(mac.Sum(nil))
}
