package internal

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Signatures represents the plugins/rules from the
// .yaml configuration file. It's the root of a config
// file.
type Signatures struct {
	Insecure string    `yaml:"insecure"`
	Plugins  []*Plugin `yaml:"plugins"`
}

// Plugin means an entry to test for during scan.
type Plugin struct {
	Endpoints       []string `yaml:"endpoints"`
	QueryString     string   `yaml:"query_string"`
	Checks          []*Check `yaml:"checks"`
	FollowRedirects bool     `yaml:"follow_redirects"`
}

// Check is a check the scan runs in.
type Check struct {
	MustMatchOne []string `yaml:"match"`
	MustMatchAll []string `yaml:"all_match"`
	MustNotMatch []string `yaml:"no_match"`
	StatusCode   *int     `yaml:"status_code"`
	Name         string   `yaml:"name"`
	Remediation  string   `yaml:"remediation"`
	Severity     string   `yaml:"severity"`
	Description  string   `yaml:"description"`
	Headers      []string `yaml:"headers"`
	NoHeaders    []string `yaml:"no_headers"`
}

// Match analyses the HTTP Response. A match means that
// one of the criteria has been met (through the strategies
// of MatchAll/MatchOne/NotMatch, and Headers/NotHeaders).
func (check *Check) Match(resp *HTTPResponse) (bool, error) {
	// Test status code
	if check.StatusCode == nil {
		return false, &ErrNilParameter{"check.StatusCode"}
	}
	if resp.StatusCode != *check.StatusCode {
		return false, nil
	}

	// Check for MatchAll
	for _, match := range check.MustMatchAll {
		if !bytes.Contains(resp.Body, []byte(match)) {
			return false, nil
		}
	}

	// Check for MatchOne
	found := false
	for _, match := range check.MustMatchOne {
		if bytes.Contains(resp.Body, []byte(match)) {
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}

	// Check for NotMatch
	for _, match := range check.MustNotMatch {
		if bytes.Contains(resp.Body, []byte(match)) {
			return false, nil
		}
	}

	// Check for headers
	for _, header := range check.Headers {
		hs := strings.Split(header, ":")
		if len(hs) != 2 {
			return false, &ErrInvalidHeaderFormat{header}
		}
		hKey := hs[0]
		hVal := hs[1]

		// Check for header in the HTTPResponse by its key
		respHVal, ok := resp.Header[hKey]
		if !ok {
			return false, nil
		}

		// Look for a match
		found = false
		for _, respHeaderValue := range respHVal {
			if strings.Contains(respHeaderValue, hVal) {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}

	// Check for NoHeaders
	for _, header := range check.NoHeaders {
		pNH := strings.Split(header, ":")
		if len(pNH) != 2 {
			return false, &ErrInvalidHeaderFormat{header}
		}

		nhKey := pNH[0]
		nhVal := pNH[1]
		if respHeaderValues, kFound := resp.Header[nhKey]; kFound {
			vFound := false
			for _, respHeaderValue := range respHeaderValues {
				if strings.Contains(respHeaderValue, nhVal) {
					vFound = true
					break
				}
			}
			if vFound {
				return false, nil
			}
		}
	}

	// If matches everything, then it's fine
	return true, nil
}

// ErrCheckInvalidField is an error meaning a check
// field is invalid.
type ErrCheckInvalidField struct {
	Check string
	Field string
}

func (e ErrCheckInvalidField) Error() string {
	return "missing or empty " + e.Field + " in " + e.Check + " plugin checks."
}

// ErrInvalidHeaderFormat is an error meaning an header
// format is invalid.
type ErrInvalidHeaderFormat struct {
	Header string
}

func (e ErrInvalidHeaderFormat) Error() string {
	return "invalid header format: " + e.Header + " should be \"KEY:VALUE\""
}

// ErrInvalidPathSignaturesFile is an error meaning
// that the path to the signatures file is invalid.
var ErrInvalidPathSignaturesFile = errors.New("path of signatures file is not valid")

// ErrBothEndpointSet is an error meaning endpoint and
// endpoints are set at same time.
var ErrBothEndpointSet = errors.New("URI and URIs can't be set at the same time in plugin checks")

// ParseSignatures parses and returns the signatures
// from the path of the file containg those.
func ParseSignatures(path string) (*Signatures, error) {
	// Check signature file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, ErrInvalidPathSignaturesFile
	}
	signFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer signFile.Close()

	// Read its content
	signData, err := io.ReadAll(signFile)
	if err != nil {
		return nil, err
	}
	var sign Signatures
	err = yaml.Unmarshal(signData, &sign)
	if err != nil {
		return nil, err
	}

	// Build signatures
	for _, plugin := range sign.Plugins {
		// Ensure the plugin's checks content are valid
		for _, check := range plugin.Checks {
			// Check main fields are not empty
			switch "" {
			case check.Description:
				return nil, &ErrCheckInvalidField{Check: check.Name, Field: "description"}
			case check.Remediation:
				return nil, &ErrCheckInvalidField{Check: check.Name, Field: "remediation"}
			case check.Severity:
				return nil, &ErrCheckInvalidField{Check: check.Name, Field: "severity"}
			}

			// Check severity is valid
			if _, err := StringToSeverity(check.Severity); err != nil {
				return nil, err
			}

			// Check headers to ensure they match KEY:VALUE fmt
			for _, header := range check.Headers {
				if strings.Count(header, ":") != 1 {
					return nil, &ErrInvalidHeaderFormat{header}
				}
			}
		}
	}

	return &sign, nil
}
