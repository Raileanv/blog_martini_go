package models

type Post struct {
	Id      int
	Title   string
	Content string
}

func NewPost(id int, title string, content string) *Post {
	return &Post{id, title, content}
}
