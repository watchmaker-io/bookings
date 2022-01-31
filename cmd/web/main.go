package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/watchmaker-io/bookings/pkg/config"
	"github.com/watchmaker-io/bookings/pkg/handlers"
	"github.com/watchmaker-io/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main entry point to app
func main() {
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.TemplateCache = tc
	app.UseCache = false
	app.Session = session

	render.NewTemplate(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
