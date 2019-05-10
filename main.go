package main

import (
	"net/http"

	"blog/models"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var posts map[string]*models.Post

func main() {
	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory: "views", // Specify what path to load the templates from.
		Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		// Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true, // Output human readable JSON
		// HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
	}))
	
	http.Handle("/assets/", http.StripPrefix("/assets/",http.FileServer(http.Dir("./assets/"))))
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
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML(200, "write", posts)
}

func editHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	post, found := posts[id]

	if !found{
		rnd.Redirect("/")
		return
	}

	rnd.HTML(200, "write", post)
}

func deleteHandler(rnd render.Render, params martini.Params) {
	id := params["id"]

	if id == "" {
		rnd.Redirect("/")
	}

	delete(posts, id)

	rnd.Redirect("/")
	
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if id != "" {
		post := posts[id]
		post.Title = title
		post.Content = content
	} else {
		id := GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	rnd.Redirect("/")
}
