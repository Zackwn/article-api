package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewForgotPasswordController(usecase *usecase.ForgotPasswordUseCase) *ForgotPasswordController {
	return &ForgotPasswordController{forgotPasswordUseCase: usecase}
}

type ForgotPasswordController struct {
	forgotPasswordUseCase *usecase.ForgotPasswordUseCase
}

func (forgotPasswordController ForgotPasswordController) Handle(req *http.Request) *Response {
	defer req.Body.Close()
	var body map[string]string
	json.NewDecoder(req.Body).Decode(&body)
	email, emailOK := body["email"]
	if !emailOK {
		return StatusResponse(http.StatusBadRequest)
	}
	dto := usecase.ForgotPasswordDTO{
		Email: email,
	}
	uErr := forgotPasswordController.forgotPasswordUseCase.Exec(dto)
	if uErr != nil {
		return ErrorResponse(uErr)
	}
	return StatusResponse(http.StatusOK)
}
