package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())

	data := struct {
		Name string
		Desc string
	}{
		"Toople",
		"gonna be super cool!",
	}

	m.Get("/", func(r render.Render) {
		r.HTML(200, "en/index", data)
	})

	m.Run()
}
