package controller

import (
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewUserProfileController(usecase *usecase.UserProfileUseCase) *UserProfileController {
	return &UserProfileController{UserProfileUseCase: usecase}
}

type UserProfileController struct {
	UserProfileUseCase *usecase.UserProfileUseCase
}

func (userProfileController UserProfileController) Handle(req *http.Request) *Response {
	id := req.URL.Query().Get("id")
	if id == "" {
		return StatusResponse(http.StatusBadRequest)
	}
	dto := usecase.UserProfileDTO{
		ID: id,
	}
	user, err := userProfileController.UserProfileUseCase.Exec(&dto)
	if err != nil {
		return ErrorResponse(err)
	}
	return JSONResponse(http.StatusOK, user)
}
