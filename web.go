package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/glacjay/goini"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/toople-co/toople-db"
)

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

func main() {
	// Load config file
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<config-file>")
		return
	}
	filename := os.Args[1]
	config, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("error reading config file '%s'\n", filename)
	}

	// Setup Database
	dba, err := setupDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Layout: "layout"}))
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("secret123"))))

	m.Get("/", func(r render.Render, s sessions.Session) {
		if id := s.Get("user"); id != nil {
			view := NewView("home", "en")
			user, err := dba.GetUser(id.(db.UserId))
			if err != nil {
				view["error"] = view["internal_error"]
				r.HTML(200, "home", view)
				return
			}
			view["loggedin_text"] = fmt.Sprintf(view["loggedin_text"].(string), user.Name())
			events, err := dba.GetEvents(id.(db.UserId))
			if err != nil {
				view["error"] = view["internal_error"]
			}
			view["events"] = events
			r.HTML(200, "home", view)
		} else {
			view := NewView("welcome", "en")
			view["login_button"] = true
			r.HTML(200, "welcome", view)
		}
	})

	m.Get("/login", func(r render.Render, s sessions.Session) {
		if s.Get("user") != nil {
			r.Redirect("/")
			return
		}
		view := NewView("login", "en")
		if email := s.Get("email"); email != nil {
			view["email"] = email
		}
		if error := s.Get("error"); error != nil {
			view["error"] = view[error.(string)]
			s.Delete("error")
		}
		r.HTML(200, "login", view)
	})

	m.Post("/login", binding.Bind(LoginForm{}), func(r render.Render, f LoginForm, s sessions.Session) {
		auth, u, err := dba.AuthUser(f.Email, f.Password)
		if auth {
			s.Set("user", u.Id())
			r.Redirect("/")
			return
		}
		switch {
		case err != nil:
			s.Set("error", "internal_error")
			fmt.Println(err)
		case u == nil:
			s.Set("error", "email_error")
		default:
			s.Set("error", "auth_error")
		}
		s.Set("email", f.Email)
		r.Redirect("/login")
	})

	m.Get("/logout", func(r render.Render, s sessions.Session) {
		s.Clear()
		r.Redirect("/")
	})

	m.Get("/signup", func(r render.Render, s sessions.Session) {
		if s.Get("user") != nil {
			r.Redirect("/")
			return
		}
		view := NewView("signup", "en")
		view["login_button"] = true
		if email := s.Get("email"); email != nil {
			view["email"] = email
		}
		if name := s.Get("name"); name != nil {
			view["name"] = name
		}
		if error := s.Get("error"); error != nil {
			view["error"] = view[error.(string)]
			s.Delete("error")
		}
		r.HTML(200, "signup", view)
	})

	m.Post("/signup", binding.Bind(SignupForm{}), func(r render.Render, f SignupForm, s sessions.Session) {
		u, err := dba.NewUser(f.Name, f.Email, f.Password)
		fmt.Println(err)
		if u != nil {
			s.Set("user", u.Id())
			r.Redirect("/")
			return
		}
		s.Set("error", "internal_error")
		s.Set("email", f.Email)
		s.Set("name", f.Name)
		r.Redirect("/signup")
	})

	m.Run()
}

func setupDatabase(config ini.Dict) (*db.DB, error) {
	dbname, ok := config.GetString("db", "name")
	if !ok {
		return nil, fmt.Errorf("config file: db name missing")
	}
	host, ok := config.GetString("db", "host")
	if !ok {
		return nil, fmt.Errorf("config file: db host missing")
	}
	port, ok := config.GetString("db", "port")
	if !ok {
		return nil, fmt.Errorf("config file: db port missing")
	}
	username, ok := config.GetString("db", "username")
	if !ok {
		return nil, fmt.Errorf("config file: db username missing")
	}
	pwd_base64, ok := config.GetString("db", "password")
	if !ok {
		return nil, fmt.Errorf("config file: db password missing")
	}
	password, err := base64.StdEncoding.DecodeString(pwd_base64)
	if err != nil {
		return nil, fmt.Errorf("config file: db password must be base64 encoded")
	}
	return db.New(host, port, username, string(password), dbname)
}
