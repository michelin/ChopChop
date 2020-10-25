package core

// Output structure for each findings
type Output struct {
	URL         string `json:"url"`
	Endpoint    string `json:"endpoint"`
	Name        string `json:"checkName"`
	Severity    string `json:"severity"`
	Remediation string `json:"remediation"`
}
