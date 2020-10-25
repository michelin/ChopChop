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

func SeverityReached(max string, sv string) bool {
	switch max {
	case "High":
		if sv == "High" {
			return true
		}
	case "Medium":
		if sv == "High" || sv == "Medium" {
			return true
		}
	case "Low":
		if sv == "High" || sv == "Medium" || sv == "Low" {
			return true
		}
	case "Informational":
		if sv == "High" || sv == "Medium" || sv == "Low" || sv == "Informational" {
			return true
		}
	}
	return false
}
