package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

func GetIndex(r render.Render, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	name := s.Get("user_name")
	if id != nil && name != nil {
		feed, err := db.GetFeed(id.(string))
		if err != nil {
			view := NewView("error", "en")
			view["error"] = view["internal_error"]
			r.HTML(500, "error", view)
			return
		}
		view := NewView("home", "en")
		view["loggedin"] = true
		view["user_name"] = name
		view["feed"] = feed
		r.HTML(200, "home", view)
		return
	}
	view := NewView("welcome", "en")
	r.HTML(200, "welcome", view)
}
