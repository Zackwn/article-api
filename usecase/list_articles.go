package usecase

import "github.com/zackwn/article-api/entity"

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

func (listArticlesUseCase ListArticlesUseCase) Exec(dto *ListArticlesDTO) ([]*entity.Article, error) {
	switch true {
	case dto.FilterOptions.ID != "":
		article, _ := listArticlesUseCase.articleRepository.FindByID(dto.FilterOptions.ID)
		return []*entity.Article{article}, nil

	case dto.FilterOptions.AuthorID != "":
		_, userFound := listArticlesUseCase.userRepository.FindByID(dto.FilterOptions.AuthorID)
		if !userFound {
			return nil, ErrUserDoNotExists{}
		}
		result := listArticlesUseCase.articleRepository.FindAllByAuthor(dto.FilterOptions.AuthorID)
		return result, nil

	// all articles
	default:
		return listArticlesUseCase.articleRepository.All(), nil
	}
}
