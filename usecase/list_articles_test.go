package usecase

import (
	"testing"

	"github.com/zackwn/article-api/entity"
)

var listArticlesUseCase *ListArticlesUseCase
var articleRepository ArticleRepository
var userRepository UserRepository

func setupListArticlesUseCase(articleRepo ArticleRepository, userRepo UserRepository) {
	articleRepository = articleRepo
	userRepository = userRepo
	listArticlesUseCase = NewListArticlesUseCase(articleRepo, userRepo)
}

type ListArticlesTestPair struct {
	filterOptions  ListArticlesFilterOptions
	NilErrorResult bool
	ResultLength   int
}

func TestListArticles(t *testing.T) {
	user, _ := entity.NewUser("name", "listarticles@mail.com", "picture", "password")
	userRepository.Store(user)
	authorID := user.ID
	article1, _ := entity.NewArticle("title article 1", "content article 1", authorID)
	article2, _ := entity.NewArticle("title article 2", "content article 2", authorID)
	article3, _ := entity.NewArticle("title article 3", "content article 3", "non existing user id")

	articleRepository.Store(article1)
	articleRepository.Store(article2)
	articleRepository.Store(article3)

	var listArticleTests []*ListArticlesTestPair = []*ListArticlesTestPair{
		{
			filterOptions:  ListArticlesFilterOptions{ID: article1.ID},
			ResultLength:   1,
			NilErrorResult: true,
		},
		{
			filterOptions:  ListArticlesFilterOptions{AuthorID: authorID},
			ResultLength:   2,
			NilErrorResult: true,
		},
		{
			filterOptions:  ListArticlesFilterOptions{},
			ResultLength:   3,
			NilErrorResult: true,
		},
		{
			filterOptions:  ListArticlesFilterOptions{ID: article1.ID, AuthorID: article1.AuthorID},
			ResultLength:   1,
			NilErrorResult: true,
		},
	}

	for _, listArticleTest := range listArticleTests {
		result, err := listArticlesUseCase.Exec(&ListArticlesDTO{FilterOptions: listArticleTest.filterOptions})
		if listArticleTest.NilErrorResult && err != nil {
			t.Errorf("Expect %v Got %v", "nil error", err)
		} else if len(result) != listArticleTest.ResultLength {
			t.Errorf("Expect %v Got %v", listArticleTest.ResultLength, len(result))
		}
	}
}

func TestListArticlesWithTwoOptions(t *testing.T) {
	user1, _ := entity.NewUser("user1", "testlistarticleswithtwoptions1@mail.com", "picture", "password123")
	user2, _ := entity.NewUser("user2", "testlistarticleswithtwoptions2@mail.com", "picture", "password123")
	userRepository.Store(user1)
	userRepository.Store(user2)
	article, _ := entity.NewArticle("article by user1", "content", user1.ID)
	articleRepository.Store(article)
	result, _ := listArticlesUseCase.Exec(&ListArticlesDTO{
		FilterOptions: ListArticlesFilterOptions{
			ID:       article.ID,
			AuthorID: user2.ID,
		},
	})
	if result[0] != nil {
		t.Errorf("Expect %v Got %v", nil, result[0])
	}
}
