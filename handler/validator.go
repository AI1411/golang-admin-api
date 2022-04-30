package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

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
	case "CreatedAt":
		return "作成日時"
	case "UpdatedAt":
		return "更新日時"
	}
	return attribute
}
