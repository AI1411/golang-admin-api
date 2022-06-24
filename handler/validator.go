package handler

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
)

const (
	validateTagForBoolean    = "boolean"
	validateTagForTime       = "datetime"
	validateTagForDateBefore = "before"
	validateTagForDateAfter  = "after"
)

var validate *validator.Validate

func init() {
	// bindする際のバリデーションルールを登録
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation(validateTagForBoolean, isBoolean)
		_ = v.RegisterValidation(validateTagForTime, isDateTimeString)
		_ = v.RegisterValidation(validateTagForDateBefore, isBeforeDateField)
		_ = v.RegisterValidation(validateTagForDateAfter, isAfterDateField)
		validate = v
	}
}

func isBoolean(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Bool:
		return true
	default:
		_, err := strconv.ParseBool(fl.Field().String())
		return err == nil
	}
}

// 日付時刻のバリデーション
func isDateTimeString(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

func isBeforeDateField(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, _, exist := fl.GetStructFieldOK2()
	if !exist || currentKind != kind {
		return false
	}

	fieldTime, err := time.Parse(time.RFC3339, field.String())
	if err != nil {
		return false
	}

	currentFieldTime, err := time.Parse(time.RFC3339, currentField.String())
	if err != nil {
		return false
	}

	return fieldTime.Before(currentFieldTime)
}

func isAfterDateField(fl validator.FieldLevel) bool {
	field := fl.Field()
	kind := field.Kind()

	currentField, currentKind, _, exist := fl.GetStructFieldOK2()
	if !exist || currentKind != kind {
		return false
	}

	fieldTime, err := time.Parse(time.RFC3339, field.String())
	if err != nil {
		return false
	}

	currentFieldTime, err := time.Parse(time.RFC3339, currentField.String())
	if err != nil {
		return false
	}

	return fieldTime.After(currentFieldTime)
}

func createValidateErrorResponse(err error) *errorResponse {
	if err == nil {
		return nil
	}
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return &errorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	details := make([]validationError, len(ve))
	for i, v := range ve {
		details[i] = validationError{
			Attribute: v.Field(),
			Message:   createValidationMessage(v.Field(), v.Tag()),
		}
	}
	return &errorResponse{
		Code:    http.StatusBadRequest,
		Message: "パラメータが不正です",
		Details: details,
	}
}

func createValidationMessage(field string, tag string) string {
	field = getAttribute(field)
	switch tag {
	case "required":
		return field + "は必須です"
	}
	return field + "は不正です"
}

func getAttribute(attribute string) string {
	switch attribute {
	case "Body":
		return "本文"
	case "Title":
		return "タイトル"
	case "UserID":
		return "ユーザーID"
	case "Status":
		return "ステータス"
	case "Age":
		return "年齢"
	case "Password":
		return "パスワード"
	case "Email":
		return "メールアドレス"
	case "PasswordConfirmation":
		return "パスワード確認"
	case "FirstName":
		return "名"
	case "LastName":
		return "姓"
	case "OrderStatus":
		return "注文ステータス"
	case "OrderID":
		return "注文ID"
	case "Remarks":
		return "備考"
	case "Quantity":
		return "数量"
	case "TotalPrice":
		return "合計金額"
	case "MilestoneTitle":
		return "タイトル"
	case "ProjectID":
		return "プロジェクトID"
	case "ProjectTitle":
		return "プロジェクト名"
	case "CreatedAt":
		return "作成日時"
	case "UpdatedAt":
		return "更新日時"
	}
	return attribute
}
