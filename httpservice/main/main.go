package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/a-h/pill/dataaccess"
	"github.com/gorilla/mux"
)

var connectionString = flag.String("connectionString", "mongodb://localhost:27017",
	"The MongoDB connection string used to store data.")

func main() {
	flag.Parse()

	r := createRoutes(dataaccess.NewMongoDataAccess(*connectionString))
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createRoutes(da dataaccess.DataAccess) *mux.Router {
	r := mux.NewRouter()

	lh := NewLoginHandler(createSession)
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

var encryptionKey = []byte("something-very-secret")

func createSession(w http.ResponseWriter, r *http.Request) Session {
	loginURL, _ := url.Parse("/")
	return NewGorillaSession(w, r, encryptionKey, *loginURL)
}
