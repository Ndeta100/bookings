package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ndeta100/booking/pkg/config"
	"github.com/Ndeta100/booking/pkg/handlers"
	"github.com/Ndeta100/booking/pkg/render"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {
	//loading environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Change this to true when in production
	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
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

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	fmt.Println(fmt.Sprintf("Staring application on port %s", os.Getenv("PORT")))
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
