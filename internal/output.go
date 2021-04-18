package internal

// Result wraps result for each finding of the scan.
type Result struct {
	URL         string `json:"url"`
	Endpoint    string `json:"endpoint"`
	Name        string `json:"checkName"`
	Severity    string `json:"severity"`
	Remediation string `json:"remediation"`
}
