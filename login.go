package main

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

func GetLogin(r render.Render, s sessions.Session) {
	if s.Get("user_id") != nil {
		r.Redirect("/")
		return
	}
	view := NewView("login", "en")
	view["email"] = s.Get("login_email")
	s.Delete("login_email")
	if error := s.Get("error"); error != nil {
		view["error"] = view[error.(string)]
		s.Delete("error")
	}
	r.HTML(200, "login", view)
}

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func PostLogin(r render.Render, f LoginForm, s sessions.Session, db *db.DB) {
	ok, u, err := db.AuthUser(f.Email, f.Password)
	if ok {
		s.Set("user_id", u.Id)
		s.Set("user_name", u.Name)
		r.Redirect("/")
		return
	}
	switch {
	case err != nil:
		s.Set("error", "internal_error")
	case u == nil:
		s.Set("error", "email_error")
	default:
		s.Set("error", "auth_error")
	}
	s.Set("login_email", f.Email)
	r.Redirect("/login")
}

func GetLogout(r render.Render, s sessions.Session) {
	s.Clear()
	r.Redirect("/")
}
