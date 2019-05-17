package main

import (
	"net/http"

	"blog/models"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	_ "github.com/lib/pq"
)

func main() {

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "views",                    // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Charset:    "UTF-8",                    // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                       // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)
	m.Post("/EditPost", editPostHandler)

	m.Run()
}

func indexHandler(rnd render.Render) {
	posts := models.AllPosts()

	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	rnd.HTML(200, "new", nil)
}

func editHandler(rnd render.Render, params martini.Params) {
	id := params["id"]

	post := models.FindPost(id)

	if post == nil {
		rnd.Redirect("/")
	}

	rnd.HTML(200, "edit", post)
}

func deleteHandler(rnd render.Render, params martini.Params) {
	id := params["id"]

	err := models.DeletePost(id)

	if err != nil {
		panic(err)
	}

	rnd.Redirect("/")
}

func editPostHandler(rnd render.Render, r *http.Request) {
	params := make(map[string]string)
	params["id"] = r.FormValue("id")
	params["title"] = r.FormValue("title")
	params["content"] = r.FormValue("content")

	err := models.UpdatePost(params)

	if err != nil {
		panic(err)
	}

	rnd.Redirect("/")
}

func savePostHandler(rnd render.Render, r *http.Request) {
	params := make(map[string]string)
	params["title"] = r.FormValue("title")
	params["content"] = r.FormValue("content")

	err := models.CreatePost(params)
	if err != nil {
		panic(err)
	}

	rnd.Redirect("/")
}
