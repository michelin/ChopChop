package core

import (
	"testing"
)

func TestValidSeverity(t *testing.T) {
	var tests = map[string]struct {
		severity string
		want     bool
	}{
		"High":          {severity: "High", want: true},
		"Medium":        {severity: "Medium", want: true},
		"Low":           {severity: "Low", want: true},
		"Informational": {severity: "Informational", want: true},
		"Bad severity":  {severity: "Unknown", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := ValidSeverity(tc.severity)
			if tc.want != have {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}

func TestSeveritiesAsString(t *testing.T) {
	want := "High, Medium, Low, Informational"
	have := SeveritiesAsString()
	if have != want {
		t.Errorf("expected: %v, got: %v", want, have)
	}
}

func TestSeverityReached(t *testing.T) {
	var tests = map[string]struct {
		max      string
		severity string
		want     bool
	}{
		"HighNotReached":       {max: "High", severity: "Informational", want: false},
		"HighReached":          {max: "High", severity: "High", want: true},
		"MediumReached":        {max: "Medium", severity: "High", want: true},
		"MediumNotReached":     {max: "Medium", severity: "Low", want: false},
		"LowReached":           {max: "Low", severity: "High", want: true},
		"LowNotReached":        {max: "Low", severity: "Informational", want: false},
		"InformationalReached": {max: "Informational", severity: "Informational", want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := SeverityReached(tc.max, tc.severity)
			if tc.want != have {
				t.Errorf("want: %v, have: %v", tc.want, have)
			}
		})
	}
}
