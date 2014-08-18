package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

func GetSignup(r render.Render, s sessions.Session) {
	if s.Get("user_id") != nil {
		r.Redirect("/")
		return
	}
	view := NewView("signup", "en")
	for _, p := range []string{"email", "name"} {
		view[p] = s.Get("signup_" + p)
		s.Delete("signup_" + p)
	}

	if error := s.Get("error"); error != nil {
		view["error"] = view[error.(string)]
		s.Delete("error")
	}
	r.HTML(200, "signup", view)
}

type SignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func PostSignup(r render.Render, f SignupForm, s sessions.Session, db *db.DB) {
	u, err := db.NewUser(f.Name, f.Email, f.Password)
	if u == nil || err != nil {
		// TODO: custom errors
		s.Set("error", "internal_error")
		s.Set("signup_email", f.Email)
		s.Set("signup_name", f.Name)
		r.Redirect("/signup")
		return
	}
	s.Set("user_id", u.Id)
	s.Set("user_name", u.Name)
	r.Redirect("/")
}
