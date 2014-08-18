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
	name := s.Get("user_name")
	view["Name"] = "Test Name"
	view["Emails"] = []string{"test@email.com", "other@hotmail.com"}
	view["loggedin"] = true
	view["user_name"] = name
	r.HTML(200, "settings", view)
}
