package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/prashant9154/Booking_System/internal/config"
	handler "github.com/prashant9154/Booking_System/internal/handlers"
	"github.com/prashant9154/Booking_System/internal/models"
	"github.com/prashant9154/Booking_System/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// what I am going to put in the session
	gob.Register(models.Reservation{})

	// change this to true in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 100 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handler.NewRepo(&app)

	handler.NewHandlers(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handler.Repo.Home)
	// http.HandleFunc("/about", handler.Repo.About)

	fmt.Printf("Starting Application on port %s \n", portNumber)

	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
