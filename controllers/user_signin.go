package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewUserSigninController(useCase *usecase.UserLoginUseCase) *UserSigninController {
	return &UserSigninController{userLoginUseCase: useCase}
}

type UserSigninDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSigninResponse struct {
	AccessToken string `json:"access_token"`
}

type UserSigninController struct {
	userLoginUseCase *usecase.UserLoginUseCase
}

func (userSigninController UserSigninController) Handle(req *http.Request) *Response {
	defer req.Body.Close()

	fmt.Println("authorization:", req.Header.Get("authorization"))

	var data UserSigninDTO
	json.NewDecoder(req.Body).Decode(&data)

	if data.Email == "" || data.Password == "" {
		return StatusResponse(http.StatusBadRequest)
	}

	accessToken, err := userSigninController.userLoginUseCase.Exec(data.Email, data.Password)
	if err != nil {
		return ErrorResponse(err)
	}

	response := UserSigninResponse{
		AccessToken: accessToken,
	}
	return JSONResponse(http.StatusOK, &response)
}
