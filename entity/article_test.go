package entity

import (
	"errors"
	"testing"
)

func TestNewArticleValid(t *testing.T) {
	_, err := NewArticle("valid title", "valid content", "valid author_id")
	if err != nil {
		t.Errorf("Expect %v Got %v", "nil", err)
	}
}

func TestNewArticleInvalid(t *testing.T) {
	_, err := NewArticle("", "valid content", "valid author_id")

	expectErr := ArticleError{"Invalid title"}
	if errors.Is(err, expectErr) == false {
		t.Errorf("Expect %v Got %v", expectErr, err)
	}
}
