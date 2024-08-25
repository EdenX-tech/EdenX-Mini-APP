package validators

import (
	"github.com/go-playground/validator/v10"
)

// NameValid 是一个自定义验证器，用于检查字符串是否为 "admin"
func NameValid(fl validator.FieldLevel) bool {
	return fl.Field().String() != "admin"
}
