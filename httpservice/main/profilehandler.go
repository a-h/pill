package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/a-h/pill/dataaccess"
)

// The ProfileHandler handles updating profile information.
type ProfileHandler struct {
	DataAccess dataaccess.DataAccess
	getSession func(w http.ResponseWriter, r *http.Request) Session
}

// NewProfileHandler creates an instance of the ProfileHandler.
func NewProfileHandler(da dataaccess.DataAccess, sessionFactory func(w http.ResponseWriter, r *http.Request) Session) *ProfileHandler {
	return &ProfileHandler{da, sessionFactory}
}

func (handler ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleProfileGet(w, r, handler)
	} else {
		handleProfilePost(w, r, handler)
	}
}

func handleProfileGet(w http.ResponseWriter, r *http.Request, handler ProfileHandler) {
	log.Printf("Handling Profile get.")

	valid, emailAddress := handler.getSession(w, r).ValidateSession()
	if !valid {
		return
	}

	log.Print("The session is valid, rendering the profile.")

	profile, _, err := handler.DataAccess.GetProfile(emailAddress)

	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve the profile for user %s.", emailAddress)
		log.Print(msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	p := &profileModel{
		Profile: profile,
		Skills:  strings.Join(flattenSkill(profile.Skills), ","),
	}

	renderTemplate(w, "profile.html", p)
}

func flattenSkill(skills []dataaccess.Skill) []string {
	op := make([]string, len(skills))

	for idx, item := range skills {
		op[idx] = item.Skill
	}

	return op
}

func handleProfilePost(w http.ResponseWriter, r *http.Request, handler ProfileHandler) {
	log.Printf("Handling Profile post.")

	valid, emailAddress := handler.getSession(w, r).ValidateSession()
	if !valid {
		return
	}

	log.Printf("The session is valid, updating the profile of %s.", emailAddress)

	// Receive post...
	err := r.ParseForm()

	if err != nil {
		log.Print("Failed to parse the form post.")
		http.Error(w, "Invalid form post.", http.StatusBadRequest)
		return
	}

	availability, _ := strconv.Atoi(r.Form.Get("availability"))

	skills := make(map[string]*dataaccess.Skill)

	for k := range r.Form {
		ok, category, group := categoriseFormKey(k)

		if !ok {
			continue
		}

		if _, ok := skills[group]; !ok {
			skills[group] = &dataaccess.Skill{}
		}

		switch category {
		case "name":
			skills[group].Skill = r.Form.Get(k)
		case "level":
			level, _ := strconv.Atoi(r.Form.Get(k))
			skills[group].Level = dataaccess.DreyfusLevel(level)
		case "interest":
			interest, _ := strconv.Atoi(r.Form.Get(k))
			skills[group].Interest = dataaccess.LikertScale(interest)
		}
	}

	pu := dataaccess.NewProfileUpdate()
	pu.Availability = dataaccess.RagStatus(availability)
	pu.EmailAddress = emailAddress
	pu.Skills = getSkillsFromMap(skills)

	_, err = handler.DataAccess.UpdateProfile(pu)

	if err != nil {
		msg := fmt.Sprintf("Unable to save profile for user %s.", emailAddress)
		log.Print(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/report/", http.StatusFound)
}

func getSkillsFromMap(skillMap map[string]*dataaccess.Skill) []dataaccess.Skill {
	skills := []dataaccess.Skill{}

	for _, v := range skillMap {
		skills = append(skills, *v)
	}

	return skills
}

// The regex used to determine whether a form post value is part of a skill.
var skillPostKeyRegex = regexp.MustCompile(`(?:name|level|interest)_(\d+)`)

func categoriseFormKey(key string) (ok bool, category string, group string) {
	if !skillPostKeyRegex.MatchString(key) {
		return false, "", ""
	}

	results := strings.Split(key, "_")
	return true, results[0], results[1]
}
