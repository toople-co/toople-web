package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

type IdForm struct {
	Id string `form:"id"`
}

func PostDismiss(r render.Render, e IdForm, s sessions.Session, db *db.DB) {
	u := s.Get("user_id")
	if u == nil {
		r.Redirect("/")
		return
	}
	err := db.DismissFeedEntry(e.Id, u.(string))
	if err != nil {
		view := NewView("error", "en")
		view["error"] = view["internal_error"]
		r.HTML(500, "error", view)
		return
	}
	r.Redirect("/")
}

func PostJoin(r render.Render, e IdForm, s sessions.Session, db *db.DB) {
	u := s.Get("user_id")
	if u == nil {
		r.Redirect("/")
		return
	}
	err := db.JoinEvent(e.Id, u.(string))
	if err != nil {
		view := NewView("error", "en")
		view["error"] = view["internal_error"]
		r.HTML(500, "error", view)
		return
	}
	r.Redirect("/")
}
