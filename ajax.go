package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

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
