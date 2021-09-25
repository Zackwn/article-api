package usecase

import "github.com/zackwn/article-api/entity"

func NewDeleteArticleUseCase(userRepo UserRepository, articleRepo ArticleRepository, auth AuthProvider) *DeleteArticleUseCase {
	return &DeleteArticleUseCase{userRepository: userRepo, articleRepository: articleRepo, authProvider: auth}
}

type DeleteArticleUseCase struct {
	userRepository    UserRepository
	articleRepository ArticleRepository
	authProvider      AuthProvider
}

type DeleteArticleDTO struct {
	AccessToken string
	ArticleID   string
}

func (deleteArticleUseCase DeleteArticleUseCase) Exec(dto *DeleteArticleDTO) UseCaseErr {
	payload, valid := deleteArticleUseCase.authProvider.Verify(dto.AccessToken)
	if !valid {
		return ErrInvalidAccessToken{}
	}
	user, found := deleteArticleUseCase.userRepository.FindByID(payload.UserID)
	if !found {
		return ErrUserDoNotExists{}
	}
	article, found := deleteArticleUseCase.articleRepository.FindByID(dto.ArticleID)
	if !found {
		return ErrArticleDoNotExists{}
	}
	if article.AuthorID != user.ID && !user.HasPermission(entity.DeleteArticlesPermission) {
		return ErrForbiddenUserAction{}
	}
	err := deleteArticleUseCase.articleRepository.Delete(dto.ArticleID)
	if err != nil {
		return ErrInternalServer{err}
	}
	return nil
}
