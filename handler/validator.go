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
	switch tag {
	case "required":
		return getAttribute(field) + "は必須です"
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
	}
	return attribute
}