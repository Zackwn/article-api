package controller

import (
	"fmt"
	"net/http"

	"github.com/zackwn/article-api/usecase"
)

func NewListArticlesController(listArticlesUseCase *usecase.ListArticlesUseCase) *ListArticlesController {
	return &ListArticlesController{listArticlesUseCase: listArticlesUseCase}
}

type ListArticlesController struct {
	listArticlesUseCase *usecase.ListArticlesUseCase
}

func (listArticlesController ListArticlesController) Handle(req *http.Request) *Response {
	queryValues := req.URL.Query()
	var authorID, ID string
	if len(queryValues["authorid"]) != 0 {
		authorID = queryValues["authorid"][0]
	}
	if len(queryValues["id"]) != 0 {
		ID = queryValues["id"][0]
	}
	dto := usecase.ListArticlesDTO{
		FilterOptions: usecase.ListArticlesFilterOptions{
			ID:       ID,
			AuthorID: authorID,
		},
	}
	articles, err := listArticlesController.listArticlesUseCase.Exec(&dto)
	fmt.Println(articles)
	if err != nil {
		return ErrorResponse(err)
	}
	return JSONResponse(http.StatusOK, articles)
}
