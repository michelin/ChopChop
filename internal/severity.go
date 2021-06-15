package internal

import (
	"strconv"
)

// Severity is a custom type defining a Signature
// severity.
// This custom type enables direct comparison between
// severities.
type Severity int

const (
	High Severity = iota
	Medium
	Low
	Informational

	highKey   = "High"
	mediumKey = "Medium"
	lowKey    = "Low"
	infoKey   = "Informational"
)

// ErrInvalidSeverity is an error meaning a given
// severity is not supported.
type ErrInvalidSeverity struct {
	Severity string
}

func (e ErrInvalidSeverity) Error() string {
	return e.Severity + " is not a valid severity"
}

// ErrUnsupportedSeverity is an error meaning a severity
// is not supported, depending on the context.
type ErrUnsupportedSeverity struct {
	Severity Severity
}

func (e ErrUnsupportedSeverity) Error() string {
	return "unsupported severity: " + strconv.Itoa(int(e.Severity))
}

func (s Severity) String() (string, error) {
	switch s {
	case High:
		return highKey, nil
	case Medium:
		return mediumKey, nil
	case Low:
		return lowKey, nil
	case Informational:
		return infoKey, nil
	default:
		return "", &ErrUnsupportedSeverity{s}
	}
}

// StringToSeverity converts a Severity to its string.
func StringToSeverity(severity string) (Severity, error) {
	switch severity {
	case highKey:
		return High, nil
	case mediumKey:
		return Medium, nil
	case lowKey:
		return Low, nil
	case infoKey:
		return Informational, nil
	default:
		return 0, &ErrInvalidSeverity{severity}
	}
}
