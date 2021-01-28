package core

import "strings"

var severities = [4]string{"High", "Medium", "Low", "Informational"}

func ValidSeverity(severity string) bool {
	for _, sv := range severities {
		if severity == sv {
			return true
		}
	}
	return false
}

func SeveritiesAsString() string {
	return strings.Join(severities[:], ", ")
}

func SeverityReached(max string, severity string) bool {
	switch max {
	case "High":
		if severity == "High" {
			return true
		}
	case "Medium":
		if severity == "High" || severity == "Medium" {
			return true
		}
	case "Low":
		if severity == "High" || severity == "Medium" || severity == "Low" {
			return true
		}
	case "Informational":
		if severity == "High" || severity == "Medium" || severity == "Low" || severity == "Informational" {
			return true
		}
	}
	return false
}
