package main

import (
	"bytes"
	"html/template"
	"log"
	"testing"

	"github.com/a-h/pill/dataaccess"
)

func TestGetLevelPcTemplateFunction(t *testing.T) {
	templateText := `{{ getlevelpc . }}`

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("test").Funcs(templateFunctions).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	tests := []struct {
		input    dataaccess.DreyfusLevel
		expected string
	}{
		{dataaccess.NoviceLevel, "80"},
		{dataaccess.CompetentLevel, "60"},
		{dataaccess.ProficientLevel, "40"},
		{dataaccess.ExpertLevel, "20"},
		{dataaccess.MasterLevel, "0"},
	}

	for _, test := range tests {
		writer := bytes.NewBufferString("")
		err = tmpl.Execute(writer, test.input)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		actual := writer.String()

		if actual != test.expected {
			t.Errorf("Expected output of '%d' to be '%s' but was '%s'",
				test.input,
				test.expected,
				actual)
		}
	}
}

func TestGetLikertPcTemplateFunction(t *testing.T) {
	templateText := `{{ getlikertpc . }}`

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("test").Funcs(templateFunctions).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	tests := []struct {
		input    dataaccess.LikertScale
		expected string
	}{
		{dataaccess.StronglyDisagree, "80"},
		{dataaccess.Disagree, "60"},
		{dataaccess.NeitherAgreeNorDisagree, "40"},
		{dataaccess.Agree, "20"},
		{dataaccess.StronglyAgree, "0"},
	}

	for _, test := range tests {
		writer := bytes.NewBufferString("")
		err = tmpl.Execute(writer, test.input)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		actual := writer.String()

		if actual != test.expected {
			t.Errorf("Expected output of '%d' to be '%s' but was '%s'",
				test.input,
				test.expected,
				actual)
		}
	}
}

func TestGetAvailabilityStyleFunction(t *testing.T) {
	templateText := `{{ getavailabilitystyle . }}`

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("test").Funcs(templateFunctions).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	tests := []struct {
		input    dataaccess.RagStatus
		expected string
	}{
		{dataaccess.Red, "alert alert-danger"},
		{dataaccess.Amber, "alert alert-warning"},
		{dataaccess.Green, "alert alert-success"},
	}

	for _, test := range tests {
		writer := bytes.NewBufferString("")
		err = tmpl.Execute(writer, test.input)
		if err != nil {
			log.Fatalf("execution: %s", err)
		}

		actual := writer.String()

		if actual != test.expected {
			t.Errorf("Expected output of '%d' to be '%s' but was '%s'",
				test.input,
				test.expected,
				actual)
		}
	}
}
