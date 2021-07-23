package usecase

import (
	"github.com/zackwn/article-api/entity"
)

func NewListArticlesUseCase(articleRepo ArticleRepository, userRepo UserRepository) *ListArticlesUseCase {
	return &ListArticlesUseCase{articleRepository: articleRepo, userRepository: userRepo}
}

type ListArticlesUseCase struct {
	articleRepository ArticleRepository
	userRepository    UserRepository
}

type ListArticlesFilterOptions struct {
	ID       string
	AuthorID string
}

type ListArticlesDTO struct {
	FilterOptions ListArticlesFilterOptions
}

func (listArticlesUseCase ListArticlesUseCase) Exec(dto *ListArticlesDTO) ([]*entity.Article, UseCaseErr) {
	// validate filter options: authorID
	if dto.FilterOptions.AuthorID != "" {
		_, userFound := listArticlesUseCase.userRepository.FindByID(dto.FilterOptions.AuthorID)
		if !userFound {
			return nil, ErrUserDoNotExists{}
		}
	}

	var articles []*entity.Article
	switch true {
	case dto.FilterOptions.AuthorID != "" && dto.FilterOptions.ID != "":
		article, _ := listArticlesUseCase.articleRepository.FindByID(dto.FilterOptions.ID)
		if article != nil && (article.AuthorID == dto.FilterOptions.AuthorID) {
			articles = []*entity.Article{article}
		} else {
			articles = []*entity.Article{nil}
		}

	case dto.FilterOptions.ID != "":
		article, _ := listArticlesUseCase.articleRepository.FindByID(dto.FilterOptions.ID)
		articles = []*entity.Article{article}

	case dto.FilterOptions.AuthorID != "":
		articles = listArticlesUseCase.articleRepository.FindAllByAuthor(dto.FilterOptions.AuthorID)

	// all articles
	default:
		articles = listArticlesUseCase.articleRepository.All()
	}
	return articles, nil
}
