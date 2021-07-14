package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewCreateArticleController(usecase *usecase.CreateArticleUseCase) *CreateArticleController {
	return &CreateArticleController{createArticleUseCase: usecase}
}

type CreateArticleController struct {
	createArticleUseCase *usecase.CreateArticleUseCase
}

func (createArticleController CreateArticleController) Handle(req *http.Request) *Response {
	defer req.Body.Close()
	accessToken := req.Header.Get("authorization")
	if accessToken == "" {
		return StatusResponse(http.StatusUnauthorized)
	}
	var result map[string]string
	json.NewDecoder(req.Body).Decode(&result)
	title, titleOK := result["title"]
	content, contentOK := result["content"]
	if !titleOK || !contentOK {
		return StatusResponse(http.StatusBadRequest)
	}
	dto := usecase.CreateArticleDTO{
		AccessToken:    accessToken,
		ArticleTitle:   title,
		ArticleContent: content,
	}
	err := createArticleController.createArticleUseCase.Exec(&dto)
	if err != nil {
		return ErrorResponse(err)
	}
	return StatusResponse(http.StatusOK)
}
