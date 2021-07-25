package controller

import (
	"fmt"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewUserSignupController(useCase *usecase.CreateUserUseCase, fileStorage FileStorage) *UserSignupController {
	return &UserSignupController{createUserUseCase: useCase, fileStorage: fileStorage}
}

type UserSignupController struct {
	createUserUseCase *usecase.CreateUserUseCase
	fileStorage       FileStorage
}

func (userSignUpController UserSignupController) Handle(req *http.Request) *Response {
	defer req.Body.Close()

	req.ParseMultipartForm(2 * 1024 * 1024)

	// verify create_user dto
	emailValue := req.MultipartForm.Value["email"]
	nameValue := req.MultipartForm.Value["name"]
	passwordValue := req.MultipartForm.Value["password"]
	if len(emailValue) != 1 || len(nameValue) != 1 || len(passwordValue) != 1 {
		return StatusResponse(http.StatusBadRequest)
	}
	email := emailValue[0]
	name := nameValue[0]
	password := passwordValue[0]
	if name == "" || password == "" || email == "" {
		return StatusResponse(http.StatusBadRequest)
	}

	// save picture
	picture := req.MultipartForm.File["picture"]
	if len(picture) != 1 {
		return StatusResponse(http.StatusBadRequest)
	}
	fileHeader := picture[0]
	file, err := fileHeader.Open()
	if err != nil {
		return StatusResponse(http.StatusBadRequest)
	}
	defer file.Close()
	pictureURL, err := userSignUpController.fileStorage.Store(file, fileHeader)
	if err != nil {
		return StatusResponse(http.StatusBadRequest)
	}

	fmt.Println(pictureURL)

	// usecase
	dto := &usecase.CreateUserDTO{
		Name:     name,
		Email:    email,
		Password: password,
		Picture:  pictureURL,
	}
	uErr := userSignUpController.createUserUseCase.Exec(dto)
	if uErr != nil {
		return ErrorResponse(uErr)
	}
	return StatusResponse(http.StatusOK)
}
