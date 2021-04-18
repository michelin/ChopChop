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
//
// XXX endpoints should only be used, not endpoint.
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
	StatusCode   *int32   `yaml:"status_code"`
	Name         string   `yaml:"name"`
	Remediation  string   `yaml:"remediation"`
	Severity     string   `yaml:"severity"`
	Description  string   `yaml:"description"`
	Headers      []string `yaml:"headers"`
	NoHeaders    []string `yaml:"no_headers"`
}

// Match analyses the HTTP Response. A match means that
// one of the criteria has been met.
//
// TODO improve this.
func (check *Check) Match(resp *HTTPResponse) bool {
	// status code must match
	if check.StatusCode != nil {
		if int32(resp.StatusCode) != *check.StatusCode {
			return false
		}
	}

	// all element must be found
	for _, match := range check.MustMatchAll {
		if !bytes.Contains(resp.Body, []byte(match)) {
			return false
		}
	}

	// one element must be found
	if len(check.MustMatchOne) > 0 {
		found := false
		for _, match := range check.MustMatchOne {
			if bytes.Contains(resp.Body, []byte(match)) {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	// no element should match
	if len(check.MustNotMatch) > 0 {
		for _, match := range check.MustNotMatch {
			if bytes.Contains(resp.Body, []byte(match)) {
				return false
			}
		}
	}

	// must contain all these headers
	for _, header := range check.Headers {
		pHeaders := strings.Split(header, ":")
		pHeadersKey := pHeaders[0]
		pHeadersValue := pHeaders[1]
		if respHeaderValues, kFound := resp.Header[pHeadersKey]; kFound {
			vFound := false
			for _, respHeaderValue := range respHeaderValues {
				if strings.Contains(respHeaderValue, pHeadersValue) {
					vFound = true
					break
				}
			}
			if !vFound {
				return false
			}
		} else {
			return false
		}
	}

	// must not contain these headers
	for _, header := range check.NoHeaders {
		pNoHeaders := strings.Split(header, ":")
		pNoHeadersKey := pNoHeaders[0]
		if respHeaderValues, kFound := resp.Header[pNoHeadersKey]; kFound {
			if len(pNoHeaders) > 1 {
				pHeadersValue := pNoHeaders[1]
				vFound := false
				for _, respHeaderValue := range respHeaderValues {
					if strings.Contains(respHeaderValue, pHeadersValue) {
						vFound = true
						break
					}
				}
				if vFound {
					return false
				}
			}
		}
	}
	return true
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
func ParseSignatures(signatures string) (*Signatures, error) {
	// Check signature file exists
	if _, err := os.Stat(signatures); os.IsNotExist(err) {
		return nil, ErrInvalidPathSignaturesFile
	}
	signFile, err := os.Open(signatures)
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
