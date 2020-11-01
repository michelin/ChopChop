package core_test

import (
	"gochopchop/core"
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
			have := core.ValidSeverity(tc.severity)
			if tc.want != have {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}

func TestSeveritiesAsString(t *testing.T) {
	want := "High, Medium, Low, Informational"
	have := core.SeveritiesAsString()
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
		"High": {max: "High", severity: "Informational", want: false},
		"Info": {max: "Informational", severity: "Informational", want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := core.SeverityReached(tc.max, tc.severity)
			if tc.want != have {
				t.Errorf("want: %v, have: %v", tc.want, have)
			}
		})
	}
}
