package main

import (
	"fmt"
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
		view[p] = s.Get("circle_" + p)
		s.Delete("circle_" + p)
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
	Circles string    `form:"circles"`
	Thresh  int       `form:"threshold"`
}

func PostEvent(r render.Render, f EventForm, s sessions.Session, db *db.DB) {
	id := s.Get("user_id")
	if id == nil {
		r.Redirect("/")
		return
	}
	fmt.Printf("\n%+v\n\n", f)
	circles := strings.Split(f.Circles, ",")
	if err := db.NewEvent(f.Date, f.Loc, f.Title, f.Info, id.(string), f.Thresh, circles); err != nil {
		fmt.Println("Date: ", f.Date)
		fmt.Println("Now : ", time.Now())
		fmt.Println(err)
		// TODO: Fix datetime
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
