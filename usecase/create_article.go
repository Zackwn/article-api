package usecase

import (
	"fmt"

	"github.com/zackwn/article-api/entity"
)

func NewCreateArticleUseCase(auth AuthProvider, articleRepo ArticleRepository) *CreateArticleUseCase {
	return &CreateArticleUseCase{authProvider: auth, articleRepository: articleRepo}
}

type CreateArticleUseCase struct {
	authProvider      AuthProvider
	articleRepository ArticleRepository
}

type CreateArticleDTO struct {
	AccessToken    string
	ArticleTitle   string
	ArticleContent string
}

func (createArticleUseCase CreateArticleUseCase) Exec(dto *CreateArticleDTO) error {
	payload, valid := createArticleUseCase.authProvider.Verify(dto.AccessToken)
	if !valid {
		return ErrInvalidAccessToken{}
	}
	fmt.Println("payload", payload)
	article, err := entity.NewArticle(dto.ArticleTitle, dto.ArticleContent, payload.UserID)
	if err != nil {
		return err
	}
	err = createArticleUseCase.articleRepository.Store(article)
	return err
}
