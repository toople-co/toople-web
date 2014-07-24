package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

func GetNewCircle(r render.Render, s sessions.Session) {
	if s.Get("user_id") == nil {
		r.Redirect("/")
		return
	}
	view := NewView("new_circle", "en")
	if name := s.Get("circle_name"); name != nil {
		view["name"] = name
		s.Delete("circle_name")
	}
	if slug := s.Get("circle_slug"); slug != nil {
		view["slug"] = slug
		s.Delete("circle_slug")
	}
	if emails := s.Get("circle_emails"); emails != nil {
		view["emails"] = emails
		s.Delete("circle_emails")
	}
	if error := s.Get("error"); error != nil {
		view["error"] = view[error.(string)]
		s.Delete("error")
	}
	r.HTML(200, "new_circle", view)
}

type CircleForm struct {
	Name   string   `form:"name"`
	Slug   string   `form:"slug"`
	Emails []string `form:"emails"`
}

func PostCircle(r render.Render, f CircleForm, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	if id == nil {
		r.Redirect("/")
		return
	}
	if _, err := db.NewCircle(f.Name, f.Slug, id.(string)); err != nil {
		s.Set("error", "internal_error")
		s.Set("circle_name", f.Name)
		s.Set("circle_slug", f.Slug)
		s.Set("circle_emails", f.Emails)
		r.Redirect("/new/circle")
		return
	}
	r.Redirect("/")
}
