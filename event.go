package main

import (
	"fmt"
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
	if date := s.Get("event_date"); date != nil {
		view["date"] = date
		s.Delete("event_date")
	}
	if title := s.Get("event_title"); title != nil {
		view["title"] = title
		s.Delete("event_title")
	}
	if loc := s.Get("event_location"); loc != nil {
		view["location"] = loc
		s.Delete("event_location")
	}
	if info := s.Get("event_info"); info != nil {
		view["info"] = info
		s.Delete("event_info")
	}
	if circles := s.Get("event_circles"); circles != nil {
		view["circles"] = circles
		s.Delete("event_circles")
	}
	if thresh := s.Get("event_threshold"); thresh != nil {
		view["threshold"] = thresh
		s.Delete("event_threshold")
	}
	if error := s.Get("error"); error != nil {
		view["error"] = view[error.(string)]
		s.Delete("error")
	}
	name := s.Get("user_name")
	view["loggedin"] = true
	view["loggedin_text"] = fmt.Sprintf(view["loggedin_text"].(string), name.(string))
	r.HTML(200, "new_event", view)
}

type EventForm struct {
	Date    time.Time `form:"date"`
	Title   string    `form:"title"`
	Loc     string    `form:"location"`
	Info    string    `form:"info"`
	Circles []string  `form:"circles"`
	Thresh  int       `form:"threshold"`
}

func PostEvent(r render.Render, f EventForm, s sessions.Session, db *db.DB) {
	id := s.Get("user")
	if id == nil {
		r.Redirect("/")
		return
	}
	if db.NewEvent(f.Date, f.Loc, f.Title, f.Info, id.(string), f.Thresh, f.Circles) != nil {
		s.Set("error", "internal_error")
		s.Set("event_date", f.Date)
		s.Set("event_title", f.Title)
		s.Set("event_location", f.Loc)
		s.Set("event_info", f.Info)
		s.Set("event_circles", f.Circles)
		s.Set("event_threshold", f.Thresh)
		r.Redirect("/event")
		return
	}
	r.Redirect("/")
}
