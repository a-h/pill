package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/a-h/pill/dataaccess"
)

// The ReportHandler handles viewing report data.
type ReportHandler struct {
	DataAccess dataaccess.DataAccess
	getSession func(w http.ResponseWriter, r *http.Request) Session
}

// NewReportHandler creates an instance of the ReportHandler.
func NewReportHandler(da dataaccess.DataAccess, sessionFactory func(w http.ResponseWriter, r *http.Request) Session) *ReportHandler {
	return &ReportHandler{da, sessionFactory}
}

func (handler ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleReportGet(w, r, handler)
}

func handleReportGet(w http.ResponseWriter, r *http.Request, handler ReportHandler) {
	log.Printf("Handling Report get.")

	valid, _ := handler.getSession(w, r).ValidateSession()
	if !valid {
		return
	}

	log.Print("The session is valid, rendering the report.")

	profiles, err := handler.DataAccess.ListProfiles()

	if err != nil {
		msg := "Unable to retrieve the list of profiles."
		log.Print(msg, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d profiles.", len(profiles))

	// List all the skills.
	skillNames := getSkillNames(profiles)
	sort.Strings(skillNames)

	log.Printf("Found %d skill names.", len(skillNames))

	profileSkills := make([]ProfileSkills, len(profiles))

	// Create a record for each profile with an array of all skills.
	for idx, profile := range profiles {
		m := &ProfileSkills{
			EmailAddress: profile.EmailAddress,
			Availability: profile.Availability,
			Skills:       make([]dataaccess.Skill, len(skillNames)),
		}

		// Create a map to speed up lookup.
		sm := make(map[string]dataaccess.Skill)
		for _, skill := range profile.Skills {
			sm[dataaccess.CleanTag(skill.Skill)] = skill
		}

		// If we have the skill, put it in the column.
		for idx, skill := range skillNames {
			value, exists := sm[skill]

			if exists {
				m.Skills[idx] = value
			}
		}

		profileSkills[idx] = *m
	}

	model := ReportModel{
		SkillNames: skillNames,
		Profiles:   profileSkills,
	}

	log.Printf("Listing %d skills.", len(model.SkillNames))
	log.Printf("Listing %d profiles.", len(model.Profiles))

	renderTemplate(w, "report.html", model)
}

// ReportModel provides data to the Report View.
type ReportModel struct {
	SkillNames []string
	Profiles   []ProfileSkills
}

// ProfileSkills provides information about a person's skills in a matrix
// format suitable for display. The ReportHandler produces a ReportModel
// which has a Skills property containing a list of all skills. The Skills
// array in this type has the same number of elements as the parent ReportModel's
// SkillNames properties, with some of the values being empty.
type ProfileSkills struct {
	EmailAddress string
	Availability dataaccess.RagStatus
	Skills       []dataaccess.Skill
}

func getSkillNames(profiles []dataaccess.Profile) []string {
	m := make(map[string]interface{})

	for _, profile := range profiles {
		for _, skill := range profile.Skills {
			m[skill.Skill] = true
		}
	}

	skillNames := make([]string, len(m))
	i := 0
	for k := range m {
		skillNames[i] = dataaccess.CleanTag(k)
		i++
	}

	return skillNames
}
