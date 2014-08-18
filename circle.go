package main

import (
	"strings"

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
	for _, p := range []string{"name", "slug", "emails"} {
		view[p] = s.Get("circle_" + p)
		s.Delete("circle_" + p)
	}
	if error := s.Get("error"); error != nil {
		view["error"] = view[error.(string)]
		s.Delete("error")
	}
	name := s.Get("user_name")
	view["loggedin"] = true
	view["user_name"] = name
	r.HTML(200, "new_circle", view)
}

type CircleForm struct {
	Name   string `form:"name"`
	Slug   string `form:"slug"`
	Emails string `form:"emails"`
}

func PostCircle(r render.Render, f CircleForm, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	if id == nil {
		r.Redirect("/")
		return
	}
	circle, err := db.NewCircle(f.Name, f.Slug, id.(string))
	if err != nil {
		s.Set("error", "internal_error")
		s.Set("circle_name", f.Name)
		s.Set("circle_slug", f.Slug)
		s.Set("circle_emails", f.Emails)
		r.Redirect("/new/circle")
		return
	}
	for _, email := range strings.Split(f.Emails, ",") {
		db.SendInvitation(circle, email)
	}

	r.Redirect("/")
}
