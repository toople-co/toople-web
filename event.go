package main

import (
	"strings"
	"time"

	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

func GetNewEvent(r render.Render, s sessions.Session) {
	if s.Get("user_id") == nil {
		r.Redirect("/")
		return
	}
	view := NewView("new_event", "en")
	for _, p := range []string{"date", "title", "location", "info", "circles", "threshold"} {
		view[p] = s.Get("event_" + p)
		s.Delete("event_" + p)
	}
	if error := s.Get("error"); error != nil {
		view["error"] = view[error.(string)]
		s.Delete("error")
	}
	name := s.Get("user_name")
	view["loggedin"] = true
	view["user_name"] = name
	r.HTML(200, "new_event", view)
}

type EventForm struct {
	Date    string `form:"date"`
	Title   string `form:"title"`
	Loc     string `form:"location"`
	Info    string `form:"info"`
	Circles string `form:"circles"`
	Thresh  int    `form:"threshold"`
}

func PostEvent(r render.Render, f EventForm, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	if id == nil {
		r.Redirect("/")
		return
	}
	// TODO: Fix datetime, refactor errors
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		s.Set("error", "internal_error")
		s.Set("event_date", f.Date)
		s.Set("event_title", f.Title)
		s.Set("event_location", f.Loc)
		s.Set("event_info", f.Info)
		s.Set("event_circles", f.Circles)
		s.Set("event_threshold", f.Thresh)
		r.Redirect("/new/event")
		return
	}
	date, err := time.ParseInLocation("Jan 2, 2006 – 15:04", f.Date, loc)
	if err != nil {
		s.Set("error", "internal_error")
		s.Set("event_date", f.Date)
		s.Set("event_title", f.Title)
		s.Set("event_location", f.Loc)
		s.Set("event_info", f.Info)
		s.Set("event_circles", f.Circles)
		s.Set("event_threshold", f.Thresh)
		r.Redirect("/new/event")
		return
	}
	circles := strings.Split(f.Circles, ",")
	if err := db.NewEvent(date, f.Loc, f.Title, f.Info, id.(string), f.Thresh, circles); err != nil {
		s.Set("error", "internal_error")
		s.Set("event_date", f.Date)
		s.Set("event_title", f.Title)
		s.Set("event_location", f.Loc)
		s.Set("event_info", f.Info)
		s.Set("event_circles", f.Circles)
		s.Set("event_threshold", f.Thresh)
		r.Redirect("/new/event")
		return
	}
	r.Redirect("/")
}
