package main

import (
	"fmt"

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
	view["loggedin"] = true
	view["loggedin_text"] = fmt.Sprintf(view["loggedin_text"].(string), name.(string))
	r.HTML(200, "settings", view)
}
