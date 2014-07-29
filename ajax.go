package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

func GetParticipants(r render.Render, p martini.Params, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	if id == nil {
		r.Redirect("/")
		return
	}
	event := p["event"]
	if event == "" {
		r.Status(404)
		return
	}
	pa, err := db.GetParticipants(event, id.(string))
	if err != nil {
		r.JSON(500, err)
		return
	}
	r.JSON(200, pa)
}

func GetCircles(r render.Render, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	if id == nil {
		r.Redirect("/")
		return
	}
	c, err := db.GetCircles(id.(string))
	if err != nil {
		r.JSON(500, err)
		return
	}
	r.JSON(200, c)
}
