package models

import (
	"database/sql"
	"log"
)

const connStr = "user=vladimir password=5b34b4ccc dbname=blog_db host=localhost sslmode=disable"

type Post struct {
	Id      int
	Title   string
	Content string
}

func NewPost(id int, title string, content string) *Post {
	return &Post{id, title, content}
}

func AllPosts() []*Post {
	var posts []*Post

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content FROM posts")
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var title, content string
	for rows.Next() {
		rows.Scan(&id, &title, &content)
		post := NewPost(id, title, content)
		posts = append(posts, post)
	}
	return posts
}

func FindPost(postId string) *Post {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatment := `
		SELECT id, title, content FROM posts
		WHERE id = $1
	`
	rows, err := db.Query(sqlStatment, postId)
	if err != nil {
		log.Fatal(err)
	}

	var post *Post
	var id int
	var title, content string
	for rows.Next() {
		rows.Scan(&id, &title, &content)
		post = NewPost(id, title, content)

	}
	return post
}

func DeletePost(postId string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
	DELETE FROM posts where id = $1
	`
	_, err = db.Exec(sqlStatement, postId)

	return err
}

func CreatePost(params map[string]string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
	INSERT INTO posts (title, content)
	VALUES ($1, $2)
	`
	_, err = db.Exec(sqlStatement, params["title"], params["content"])

	return err
}

func UpdatePost(params map[string]string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
UPDATE posts
SET title = $2, content = $3
WHERE id = $1;
`
	_, err = db.Exec(sqlStatement, params["id"], params["title"], params["content"])

	return err
}
