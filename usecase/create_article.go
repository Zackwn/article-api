package usecase

import (
	"github.com/zackwn/article-api/entity"
)

func NewCreateArticleUseCase(auth AuthProvider, articleRepo ArticleRepository, userRepo UserRepository) *CreateArticleUseCase {
	return &CreateArticleUseCase{authProvider: auth, articleRepository: articleRepo, userRepository: userRepo}
}

type CreateArticleUseCase struct {
	authProvider      AuthProvider
	articleRepository ArticleRepository
	userRepository    UserRepository
}

type CreateArticleDTO struct {
	AccessToken    string
	ArticleTitle   string
	ArticleContent string
}

func (createArticleUseCase CreateArticleUseCase) Exec(dto *CreateArticleDTO) UseCaseErr {
	// check access token
	payload, valid := createArticleUseCase.authProvider.Verify(dto.AccessToken)
	if !valid {
		return ErrInvalidAccessToken{}
	}
	user, found := createArticleUseCase.userRepository.FindByID(payload.UserID)
	if !found {
		return ErrUserDoNotExists{}
	}
	// check user permission
	if !user.HasPermission(entity.CreateArticlesPermission) {
		return ErrForbiddenUserAction{}
	}
	// create article
	article, err := entity.NewArticle(dto.ArticleTitle, dto.ArticleContent, user.ID)
	if err != nil {
		return ErrInvalidEntityField{errMessage: err.Error()}
	}
	err = createArticleUseCase.articleRepository.Store(article)
	if err != nil {
		return ErrInvalidEntityField{errMessage: err.Error()}
	}
	return nil
}
