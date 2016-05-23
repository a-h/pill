package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/a-h/pill/dataaccess"
	"github.com/a-h/pill/tokenverifier"
	"github.com/gorilla/mux"
)

var connectionString = flag.String("connectionString", "mongodb://mongo:27017",
	"The MongoDB connection string used to store data.")

func main() {
	log.Print("Starting up...")
	flag.Parse()

	log.Print("Connecting to MongoDB to retrieve configuration.")
	da := dataaccess.NewMongoDataAccess(*connectionString, "pill")

	var err error
	configuration, err = da.GetOrCreateConfiguration()

	if err != nil {
		log.Fatal("Failed to retrieve configuration, the application cannot start. ", err)
	}

	log.Print("Configuration retrieved.")

	log.Print("Creating routes...")
	r := createRoutes(da)

	log.Print("Serving...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createRoutes(da dataaccess.DataAccess) *mux.Router {
	r := mux.NewRouter()

	lh := NewLoginHandler(createSession, tokenverifier.GoogleTokenVerifier{})
	r.Handle("/", lh)

	ph := NewProfileHandler(da, createSession)
	r.Handle("/profile/", ph)

	sh := NewSkillHandler(da, createSession)
	r.Handle("/skills/", sh)

	rh := NewReportHandler(da, createSession)
	r.Handle("/report/", rh)

	// Serve static content.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return r
}

var configuration dataaccess.Configuration

func createSession(w http.ResponseWriter, r *http.Request) Session {
	loginURL, _ := url.Parse("/")
	return NewGorillaSession(w, r, configuration.SessionEncryptionKey, configuration.SetSecureFlag, *loginURL)
}
