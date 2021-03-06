package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/a-h/pill/dataaccess"
)

var templates = template.Must(template.New("").Funcs(templateFunctions).ParseFiles("templates/header.html",
	"./templates/navigation.html",
	"./templates/login.html",
	"./templates/profile.html",
	"./templates/report.html",
	"./templates/footer.html"))

// Register the functions required to render the templates.
var templateFunctions = template.FuncMap{
	"getlevelpc":           getlevelpc,
	"getlikertpc":          getlikertpc,
	"getavailabilitystyle": getavailabilitystyle,
}

// Helper function to render templates.
func renderTemplate(w http.ResponseWriter, templateName string, model interface{}) {
	err := templates.ExecuteTemplate(w, templateName, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Template functions.
func getlevelpc(a dataaccess.DreyfusLevel) string {
	return strconv.FormatInt(100-(int64(a)*20), 10)
}

func getlikertpc(a dataaccess.LikertScale) string {
	return strconv.FormatInt(100-(int64(a)*20), 10)
}

func getavailabilitystyle(a dataaccess.RagStatus) string {
	switch a {
	case dataaccess.Red:
		return "alert alert-danger"
	case dataaccess.Amber:
		return "alert alert-warning"
	}

	return "alert alert-success"
}
