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
	db, err := setupDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	m := martini.Classic()

	m.Map(db)
	m.Use(render.Renderer(render.Options{Layout: "layout"}))
	m.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("secret123"))))

	// Welcome page / home page
	m.Get("/", GetIndex)

	// Login / logout
	m.Get("/login", GetLogin)
	m.Post("/login", binding.Bind(LoginForm{}), PostLogin)
	m.Get("/logout", GetLogout)

	// Signup
	m.Get("/signup", GetSignup)
	m.Post("/signup", binding.Bind(SignupForm{}), PostSignup)

	// New event page
	m.Get("/new/event", GetNewEvent)
	m.Post("/event", binding.Bind(EventForm{}), PostEvent)
	// TODO: m.Post("/join/event", binding.Bind(JoinEventForm{}), PostJoinEvent)

	// New circle page
	m.Get("/new/circle", GetNewCircle)
	m.Post("/circle", binding.Bind(CircleForm{}), PostCircle)

	// Settings (account, profile)
	m.Get("/settings", GetSettings)

	// Ajax api
	m.Get("/participants/:event", GetParticipants)
	m.Get("/circles", GetCircles)

	// 404: Not found
	m.NotFound(GetNotFound)

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
