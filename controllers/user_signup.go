package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewUserSignupController(useCase *usecase.CreateUserUseCase) *UserSignupController {
	return &UserSignupController{createUserUseCase: useCase}
}

type UserSignupDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserSignupController struct {
	createUserUseCase *usecase.CreateUserUseCase
}

func (userSignUpController UserSignupController) Handle(req *http.Request) *Response {
	defer req.Body.Close()

	var data UserSignupDTO
	json.NewDecoder(req.Body).Decode(&data)

	if data.Name == "" || data.Password == "" || data.Email == "" {
		return StatusResponse(http.StatusBadRequest)
	}

	// usecase
	err := userSignUpController.createUserUseCase.Exec(data.Name, data.Email, data.Password)
	if err != nil {
		return ErrorResponse(err)
	}
	return StatusResponse(http.StatusOK)
}
