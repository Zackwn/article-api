package controller

import (
	"encoding/json"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewDeleteArticleController(usecase *usecase.DeleteArticleUseCase) *DeleteArticleController {
	return &DeleteArticleController{deleteArticleUseCase: usecase}
}

type DeleteArticleController struct {
	deleteArticleUseCase *usecase.DeleteArticleUseCase
}

func (deleteArticleController DeleteArticleController) Handle(req *http.Request) *Response {
	accessToken := req.Header.Get("authorization")
	var result map[string]string
	json.NewDecoder(req.Body).Decode(&result)
	articleid, articleidOK := result["id"]
	if !articleidOK {
		return StatusResponse(http.StatusBadRequest)
	}
	dto := usecase.DeleteArticleDTO{
		AccessToken: accessToken,
		ArticleID:   articleid,
	}
	err := deleteArticleController.deleteArticleUseCase.Exec(&dto)
	if err != nil {
		return ErrorResponse(err)
	}
	return StatusResponse(http.StatusOK)
}
