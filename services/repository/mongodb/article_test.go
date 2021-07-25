package mongodb

import (
	"log"
	"testing"

	"github.com/zackwn/article-api/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

var articleRepo *ArticleRepository

func setupArticleTests(db *mongo.Database) {
	articleRepo = NewArticleRepository(db)
}

func TestArticleFindAllByAuthor(t *testing.T) {
	authorID := "bd04199c-512c-4cbe-a32c-863fc827c08f"
	article1, _ := entity.NewArticle("article 1", "content", authorID)
	article2, _ := entity.NewArticle("article 2", "content", authorID)
	err := articleRepo.Store(article1)
	if err != nil {
		log.Fatal(err)
	}
	err = articleRepo.Store(article2)
	if err != nil {
		log.Fatal(err)
	}
	articles := articleRepo.FindAllByAuthor(authorID)
	if len(articles) != 2 {
		t.Errorf("Expect %v Got %v", "articles length to be 2", len(articles))
	}
}

func TestArticleFindByIDExisting(t *testing.T) {
	article, _ := entity.NewArticle("title", "content", "b21919de-06e5-4fb0-a89d-f7502377d1c7")
	_ = articleRepo.Store(article)
	_, found := articleRepo.FindByID(article.ID)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
}

func TestArticleFindByIDNotExisting(t *testing.T) {
	_, found := articleRepo.FindByID("9bfad6fe-2ad5-4c8d-a060-e71a5a98db47")
	if found != false {
		t.Errorf("Expect %v Got %v", "false", found)
	}
}

func TestArticlesAll(t *testing.T) {
	article, _ := entity.NewArticle("title", "content", "b29f69a3-883c-446b-a8e7-bc1ca63435a6")
	articleRepo.Store(article)
	articles := articleRepo.All()
	if len(articles) < 1 {
		t.Errorf("Expect %v Got %v", "articles length to be more than 0", len(articles))
	}
}
