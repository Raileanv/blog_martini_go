package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"blog/models"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"database/sql"

	_ "github.com/lib/pq"
)

var posts map[int]*models.Post

func main() {

	posts = make(map[int]*models.Post, 0)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "views",                    // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Charset:    "UTF-8",                    // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                       // Output human readable JSON
	}))

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	m.Run()
}

func indexHandler(rnd render.Render) {
	connStr := "user=vrailean password=******** dbname=blog_db host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content FROM posts")
	var id int
	var title, content string
	for rows.Next() {
		rows.Scan(&id, &title, &content)
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML(200, "new", nil)
}

func editHandler(rnd render.Render, params martini.Params) {
	id := params["id"]

	convertedID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}
	post, found := posts[convertedID]

	if !found {
		rnd.Redirect("/")
		return
	}

	rnd.HTML(200, "edit", post)
}

func deleteHandler(rnd render.Render, params martini.Params) {
	connStr := "user=vrailean password=5B34b4ddcc dbname=blog_db host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := params["id"]

	if id == "" {
		rnd.Redirect("/")
	}

	convertedID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
	}

	delete(posts, convertedID)
	sqlStatement := `
	DELETE FROM posts where id = $1
	`
	_, err = db.Exec(sqlStatement, convertedID)
	if err != nil {
		panic(err)
	}

	rnd.Redirect("/")

}

func savePostHandler(rnd render.Render, r *http.Request) {
	connStr := "user=vrailean password=5B34b4ddcc dbname=blog_db host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if id != "" {

		convertedID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
		}
		sqlStatement := `
	UPDATE posts
	SET title = $2, content = $3
	WHERE id = $1;
	`
		_, err = db.Exec(sqlStatement, convertedID, title, content)
		if err != nil {
			panic(err)
		}

	} else {
		sqlStatement := `
		INSERT INTO posts (title, content) 
		VALUES ($1, $2)
		`

		_, err = db.Exec(sqlStatement, title, content)
		if err != nil {
			panic(err)
		}
	}

	rnd.Redirect("/")
}
