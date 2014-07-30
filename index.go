package main

import (
	"fmt"

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
		view["loggedin_text"] = fmt.Sprintf(view["loggedin_text"].(string), name.(string))
		view["feed"] = feed
		r.HTML(200, "home", view)
		return
	}
	view := NewView("welcome", "en")
	view["login_button"] = true
	r.HTML(200, "welcome", view)
}

// Feed: (sorted by date)
// ====
//
//   Event (date > now && # < thresh) => Pending
//     Participant
//       User
//     Participant
//       User
//     ...
//   Event (date < now && # < thresh) => No Toople :(
//   Event (date > now && # >= thresh) => Toople :)
//     Participant
//       User
//     ...
//   Member (not me) => Alert: someone joined
//   Member (me) => You're in
//
