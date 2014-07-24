package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

func GetSettings(r render.Render, s sessions.Session) {
	if s.Get("user_id") == nil {
		r.Redirect("/")
		return
	}
	view := NewView("settings", "en")
	r.HTML(200, "settings", view)
}
