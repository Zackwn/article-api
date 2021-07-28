package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewChangePasswordController(usecase *usecase.ChangePasswordUseCase) *ChangePasswordController {
	return &ChangePasswordController{changePasswordUseCase: usecase}
}

type ChangePasswordController struct {
	changePasswordUseCase *usecase.ChangePasswordUseCase
}

func (changePasswordController ChangePasswordController) Handle(req *http.Request) *Response {
	defer req.Body.Close()
	var body map[string]string
	json.NewDecoder(req.Body).Decode(&body)
	newPassword, newPassOK := body["new_password"]
	token, tokenOK := body["token"]
	if !newPassOK || !tokenOK {
		return StatusResponse(http.StatusBadRequest)
	}
	dto := &usecase.ChangePasswordDTO{
		Token:       token,
		NewPassword: newPassword,
	}
	uErr := changePasswordController.changePasswordUseCase.Exec(dto)
	if uErr != nil {
		return ErrorResponse(uErr)
	}
	return StatusResponse(http.StatusOK)
}
