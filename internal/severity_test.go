package internal_test

import (
	"testing"

	"github.com/michelin/gochopchop/internal"
)

func TestSeverityString(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Severity    internal.Severity
		ExpectedStr string
		ExpectedErr error
	}{
		"high": {
			Severity:    internal.High,
			ExpectedStr: "High",
			ExpectedErr: nil,
		},
		"medium": {
			Severity:    internal.Medium,
			ExpectedStr: "Medium",
			ExpectedErr: nil,
		},
		"low": {
			Severity:    internal.Low,
			ExpectedStr: "Low",
			ExpectedErr: nil,
		},
		"info": {
			Severity:    internal.Informational,
			ExpectedStr: "Informational",
			ExpectedErr: nil,
		},
		"invalid": {
			Severity:    -1,
			ExpectedStr: "",
			ExpectedErr: &internal.ErrUnsupportedSeverity{-1},
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			str, err := tt.Severity.String()

			if str != tt.ExpectedStr {
				t.Error("Failed to get expected Severity.String() result: got \"" + str + "\" instead of \"" + tt.ExpectedStr)
			}

			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

func TestStringToSeverity(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Severity         string
		ExpectedSeverity internal.Severity
		ExpectedErr      error
	}{
		"high": {
			Severity:         "High",
			ExpectedSeverity: internal.High,
			ExpectedErr:      nil,
		},
		"medium": {
			Severity:         "Medium",
			ExpectedSeverity: internal.Medium,
			ExpectedErr:      nil,
		},
		"low": {
			Severity:         "Low",
			ExpectedSeverity: internal.Low,
			ExpectedErr:      nil,
		},
		"info": {
			Severity:         "Informational",
			ExpectedSeverity: internal.Informational,
			ExpectedErr:      nil,
		},
		"invalid": {
			Severity:         "invalid",
			ExpectedSeverity: 0,
			ExpectedErr:      &internal.ErrInvalidSeverity{"invalid"},
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			sev, err := internal.StringToSeverity(tt.Severity)

			if sev != tt.ExpectedSeverity {
				t.Error("Failed to get expected severity: got \"", sev, "\" instead of \"", tt.ExpectedSeverity, "\"")
			}

			checkErr(err, tt.ExpectedErr, t)
		})
	}
}
