package main

import "github.com/martini-contrib/render"

func GetNotFound(r render.Render) {
	view := NewView("error", "en")
	r.HTML(200, "error", view)
}
