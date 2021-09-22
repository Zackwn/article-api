package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewVerifyAccountController(usecase *usecase.VerifyAccountUseCase) *VerifyAccountController {
	return &VerifyAccountController{VerifyAccountUseCase: usecase}
}

type VerifyAccountController struct {
	VerifyAccountUseCase *usecase.VerifyAccountUseCase
}

func (verifyAccountController VerifyAccountController) Handle(req *http.Request) *Response {
	defer req.Body.Close()
	var result map[string]string
	json.NewDecoder(req.Body).Decode(&result)
	token, tokenOK := result["token"]
	if !tokenOK {
		return StatusResponse(http.StatusBadRequest)
	}
	dto := usecase.VerifyAccountDTO{
		Token: token,
	}
	err := verifyAccountController.VerifyAccountUseCase.Exec(&dto)
	if err != nil {
		return ErrorResponse(err)
	}
	return StatusResponse(http.StatusOK)
}
