package internal

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

type ErrInvalidSeverity struct {
	severity string
}

func (e ErrInvalidSeverity) Error() string {
	return e.severity + " is not a valid severity"
}

func (s Severity) String() string {
	switch s {
	case High:
		return highKey
	case Medium:
		return mediumKey
	case Low:
		return lowKey
	case Informational:
		return infoKey
	default:
		return "UNSUPPORTED SEVERITY"
	}
}

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
