package internal

// Output wraps result for each finding
// TODO rename
type Output struct {
	URL         string `json:"url"`
	Endpoint    string `json:"endpoint"`
	Name        string `json:"checkName"`
	Severity    string `json:"severity"`
	Remediation string `json:"remediation"`
}
