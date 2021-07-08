package controller

import (
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewTestAuthController(usecase *usecase.TestAuthUseCase) *TestAuthController {
	return &TestAuthController{testAuthUseCase: usecase}
}

type TestAuthController struct {
	testAuthUseCase *usecase.TestAuthUseCase
}

func (testAuthController TestAuthController) Handle(req *http.Request) *Response {
	defer req.Body.Close()
	accessToken := req.Header.Get("authorization")
	if accessToken == "" {
		return StatusResponse(http.StatusUnauthorized)
	}
	testAuthController.testAuthUseCase.Exec(accessToken)
	return StatusResponse(http.StatusOK)
}
