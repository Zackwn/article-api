package entity

type Article struct {
	ID        string `json:"id"`
	Publisher *User
	Title     string `json:"title"`
	Content   string `json:"content"`
}
