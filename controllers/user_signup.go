package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewUserSignupController(CreateUserUseCase *usecase.CreateUserUseCase, fileStorage FileStorage, SVATUseCase *usecase.SendVerifyAccountTokenUseCase) *UserSignupController {
	return &UserSignupController{createUserUseCase: CreateUserUseCase, fileStorage: fileStorage, sendVerifyAccountToken: SVATUseCase}
}

type UserSignupController struct {
	createUserUseCase      *usecase.CreateUserUseCase
	sendVerifyAccountToken *usecase.SendVerifyAccountTokenUseCase
	fileStorage            FileStorage
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
	pictureURL, filename, err := userSignUpController.fileStorage.Store(file, fileHeader)
	if err != nil {
		return StatusResponse(http.StatusBadRequest)
	}

	fmt.Println(pictureURL)

	// create user usecase
	dto := &usecase.CreateUserDTO{
		Name:     name,
		Email:    email,
		Password: password,
		Picture:  pictureURL,
	}
	user, uErr := userSignUpController.createUserUseCase.Exec(dto)
	if uErr != nil {
		err = userSignUpController.fileStorage.Discard(filename)
		if err != nil {
			log.Println(err)
		}
		return ErrorResponse(uErr)
	}
	// send verify account token usecase
	svatDto := &usecase.SendVerifyAccountTokenDTO{
		UserID: user.ID,
	}
	uErr = userSignUpController.sendVerifyAccountToken.Exec(svatDto)
	if uErr != nil {
		return ErrorResponse(uErr)
	}
	return StatusResponse(http.StatusOK)
}
