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
		{dataaccess.NoviceLevel, "20"},
		{dataaccess.CompetentLevel, "40"},
		{dataaccess.ProficientLevel, "60"},
		{dataaccess.ExpertLevel, "80"},
		{dataaccess.MasterLevel, "100"},
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
