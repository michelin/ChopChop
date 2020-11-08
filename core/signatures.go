package core

import (
	"gochopchop/internal"
	"strings"
)

// Signature struct to load the plugins/rules from the YAML file
type Signatures struct {
	Plugins []*Plugin `yaml:"plugins"`
}

type Plugin struct {
	Endpoints       []string `yaml:"endpoints"`
	Endpoint        string   `yaml:"endpoint"`
	QueryString     string   `yaml:"query_string"`
	Checks          []*Check `yaml:"checks"`
	FollowRedirects bool     `yaml:"follow_redirects"`
}

// Check Signature
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

// NewSignatures returns a new initialized Signatures
func NewSignatures() *Signatures {
	return &Signatures{}
}

func (s *Signatures) FilterBySeverity(severity string) {
	filteredPlugins := s.Plugins[:0]
	for _, plugin := range s.Plugins {
		filteredChecks := plugin.Checks[:0]
		for _, check := range plugin.Checks {
			if check.Severity == severity {
				filteredChecks = append(filteredChecks, check)
			}
		}
		if len(filteredChecks) > 0 {
			plugin.Checks = filteredChecks
			filteredPlugins = append(filteredPlugins, plugin)
		}
	}
	s.Plugins = filteredPlugins
}

func (s *Signatures) FilterByNames(names []string) {
	filteredPlugins := s.Plugins[:0]
	for _, plugin := range s.Plugins {
		filteredChecks := plugin.Checks[:0]
		for _, check := range plugin.Checks {
			for _, name := range names {
				if strings.Contains(strings.ToLower(check.Name), strings.ToLower(name)) {
					filteredChecks = append(filteredChecks, check)
					break
				}
			}
		}
		if len(filteredChecks) > 0 {
			plugin.Checks = filteredChecks
			filteredPlugins = append(filteredPlugins, plugin)
		}
	}
	s.Plugins = filteredPlugins
}

//Match analyses the HTTP Request
func (check *Check) Match(resp *internal.HTTPResponse) bool {
	// status code must match
	if check.StatusCode != nil {
		if int32(resp.StatusCode) != *check.StatusCode {
			return false
		}
	}
	// all element must be found
	for _, match := range check.MustMatchAll {
		if !strings.Contains(resp.Body, match) {
			return false
		}
	}

	// one element must be found
	if len(check.MustMatchOne) > 0 {
		found := false
		for _, match := range check.MustMatchOne {
			if strings.Contains(resp.Body, match) {
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
			if strings.Contains(resp.Body, match) {
				return false
			}
		}
	}

	// must contain all these headers
	if len(check.Headers) > 0 {
		for _, header := range check.Headers {
			pHeaders := strings.Split(header, ":")
			if v, kFound := resp.Header[pHeaders[0]]; kFound {
				vFound := false
				for _, n := range v {
					if strings.Contains(n, pHeaders[1]) {
						vFound = true
					}
				}
				if !vFound {
					return false
				}
			} else {
				return false
			}
		}
	}

	// must not contain these headers
	if len(check.NoHeaders) > 0 {
		for _, header := range check.NoHeaders {
			pNoHeaders := strings.Split(header, ":")
			if v, kFound := resp.Header[pNoHeaders[0]]; kFound {
				return false
			} else if kFound && len(pNoHeaders) == 1 { // if the header has not been specified.
				return false
			} else {
				for _, n := range v {
					if strings.Contains(n, pNoHeaders[1]) {
						return false
					}
				}
			}
		}
	}

	return true
}

func (self *Signatures) Equals(signatures *Signatures) bool {
	if len(self.Plugins) != len(signatures.Plugins) {
		return false
	}
	for _, plugin := range self.Plugins {
		found := false
		for _, oplugin := range signatures.Plugins {
			if plugin.Equals(oplugin) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (self *Plugin) Equals(plugin *Plugin) bool {
	if !sliceStringEqual(self.Endpoints, plugin.Endpoints) {
		return false
	}
	if self.Endpoint != plugin.Endpoint {
		return false
	}
	if self.QueryString != plugin.QueryString {
		return false
	}
	if self.FollowRedirects != plugin.FollowRedirects {
		return false
	}
	for _, check := range self.Checks {
		found := false
		for _, pcheck := range plugin.Checks {
			if check.Equals(pcheck) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (self *Check) Equals(check *Check) bool {
	if !sliceStringEqual(self.MustMatchOne, check.MustMatchOne) {
		return false
	}
	if !sliceStringEqual(self.MustMatchAll, check.MustMatchAll) {
		return false
	}
	if !sliceStringEqual(self.MustNotMatch, check.MustNotMatch) {
		return false
	}
	if self.StatusCode != nil && check.StatusCode != nil {
		if *self.StatusCode != *check.StatusCode {
			return false
		}
	}
	if self.Name != check.Name {
		return false
	}
	if self.Remediation != check.Remediation {
		return false
	}
	if self.Severity != check.Severity {
		return false
	}
	if self.Description != check.Description {
		return false
	}
	if !sliceStringEqual(self.Headers, check.Headers) {
		return false
	}
	if !sliceStringEqual(self.NoHeaders, check.NoHeaders) {
		return false
	}
	return true
}

func sliceStringEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
