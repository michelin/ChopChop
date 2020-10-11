package data

import "errors"

// Output structure for each findings
type Output struct {
	Domain      string
	TestedURL   string
	PluginName  string
	Severity    string
	Remediation string
}

// SeverityType is basically an enum and values can be from Info, Low, Medium and High
type SeverityType string

const (
	// Informational will be the default severityType
	Informational SeverityType = "Informational"
	// Low severity
	Low = "Low"
	// Medium severity
	Medium = "Medium"
	// High severity (highest rating)
	High = "High"
)

// Check Signature
type Check struct {
	Match       []*string     `yaml:"match"`
	AllMatch    []*string     `yaml:"all_match"`
	StatusCode  *int32        `yaml:"status_code"`
	PluginName  string        `yaml:"name"`
	Remediation *string       `yaml:"remediation"`
	Severity    *SeverityType `yaml:"severity"`
	Description *string       `yaml:"description"`
	NoMatch     []*string     `yaml:"no_match"`
	Headers     []*string     `yaml:"headers"`
}

// Config struct to load the configuration from the YAML file
type Config struct {
	Insecure bool        `yaml:"insecure"`
	Plugins  []Signature `yaml:"plugins"`
}

// Signature struct to load it afterwards
type Signature struct {
	URI             string  `yaml:"uri"`
	QueryString     string  `yaml:"query_string"`
	Checks          []Check `yaml:"checks"`
	FollowRedirects *bool   `yaml:"follow_redirects"`
}

// IsValid will verify that the severityType is part of the enum previously declared
func (st SeverityType) IsValid() error {
	switch st {
	case Informational, Low, Medium, High:
		return nil
	}
	return errors.New("Invalid Severity type. Please Check yaml config file")
}

func (st SeverityType) String() string {
	return string(st)
}
