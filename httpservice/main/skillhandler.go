package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/a-h/pill/dataaccess"
)

// SkillHandler lists all of the skills previously mentioned.
type SkillHandler struct {
	DataAccess dataaccess.DataAccess
	getSession func(w http.ResponseWriter, r *http.Request) Session
}

// NewSkillHandler creates an instance of the SkillHandler.
func NewSkillHandler(da dataaccess.DataAccess, sessionFactory func(w http.ResponseWriter, r *http.Request) Session) *SkillHandler {
	return &SkillHandler{da, sessionFactory}
}

func (sh SkillHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling skill request.")

	skillTags, err := sh.DataAccess.ListSkillTags()

	if err != nil {
		log.Printf("Failed to list skill tags, with error %s", err)
		http.Error(w, "Failed to list skill tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(skillTags); err != nil {
		log.Printf("Failed to marshall the skill tags, with error %s", err)
		http.Error(w, "Failed to marshall skill tags", http.StatusInternalServerError)
	}
}
