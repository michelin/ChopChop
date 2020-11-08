package core

import (
	"testing"
)

func createInt32(x int32) *int32 {
	return &x
}

// Checks

var FakeCheckStatusCode200 = &Check{
	Name:        "StatusCode200",
	Severity:    "Medium",
	Remediation: "uninstall",
	StatusCode:  createInt32(200),
}

var FakeCheckStatusCode500 = &Check{
	Name:        "StatusCode500",
	Severity:    "High",
	Remediation: "uninstall",
	StatusCode:  createInt32(500),
}

var FakeCheckNoHeaders = &Check{
	Name:        "NoHeaders",
	Severity:    "Medium",
	Remediation: "uninstall",
	NoHeaders:   []string{"NoHeader:ok"},
}

var FakeCheckHeaders = &Check{
	Name:        "Headers",
	Severity:    "Medium",
	Remediation: "uninstall",
	Headers:     []string{"Header:ok"},
}

var FakeCheckMatchOne = &Check{
	Name:         "MustMatchOne",
	Severity:     "Medium",
	Remediation:  "uninstall",
	MustMatchOne: []string{"MATCHONE", "MATCHTWO"},
}

var FakeCheckMatchAll = &Check{
	Name:         "MustMatchAll",
	Severity:     "Medium",
	Remediation:  "uninstall",
	MustMatchAll: []string{"MATCHONE", "MATCHTWO"},
}

var FakeCheckNotMatch = &Check{
	Name:         "MustNotMatch",
	Severity:     "Medium",
	Remediation:  "uninstall",
	MustNotMatch: []string{"NOTMATCH"},
}

// Plugins

var FakePlugin = &Plugin{
	Endpoint: "/",
	Checks: []*Check{
		FakeCheckStatusCode200,
		FakeCheckHeaders,
		FakeCheckNoHeaders,
		FakeCheckMatchAll,
		FakeCheckMatchOne,
		FakeCheckNotMatch,
	},
}

var FakeQueryPlugin = &Plugin{
	Endpoint:    "/",
	QueryString: "query=test",
	Checks: []*Check{
		FakeCheckStatusCode200,
	},
}

var FakePlugin2 = &Plugin{
	Endpoint:    "/fake",
	QueryString: "query=test",
	Checks: []*Check{
		FakeCheckStatusCode500,
	},
}

var FakeFollowRedirectPlugin = &Plugin{
	Endpoint: "/",
	Checks: []*Check{
		FakeCheckStatusCode200,
	},
	FollowRedirects: true,
}

// Signatures
var FakeSignatures = &Signatures{
	Plugins: []*Plugin{
		FakePlugin,
		FakeQueryPlugin,
		FakeFollowRedirectPlugin,
	},
}

func TestFilterBySeverity(t *testing.T) {
	var tests = map[string]struct {
		have     *Signatures
		want     *Signatures
		severity string
	}{
		"Filter nothing": {
			have:     &Signatures{Plugins: []*Plugin{FakePlugin}},
			want:     &Signatures{Plugins: []*Plugin{FakePlugin}},
			severity: "Medium",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.have.FilterBySeverity(tc.severity)
			if !tc.want.Equals(tc.have) {
				t.Errorf("expected: %v, got: %v", tc.want, tc.have)
			}
		})
	}
}

func TestFilterByNames(t *testing.T) {
	var tests = map[string]struct {
		have  *Signatures
		want  *Signatures
		names []string
	}{
		"Filter nothing": {
			have:  &Signatures{Plugins: []*Plugin{FakeQueryPlugin}},
			want:  &Signatures{Plugins: []*Plugin{FakeQueryPlugin}},
			names: []string{FakeCheckStatusCode200.Name},
		},
		"Filter one element": {
			have:  &Signatures{Plugins: []*Plugin{FakeQueryPlugin}},
			want:  &Signatures{},
			names: []string{"check's name that is not in the signatures"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.have.FilterByNames(tc.names)
			if !tc.have.Equals(tc.want) {
				t.Errorf("expected: %v, got: %v", tc.want, tc.have)
			}
		})
	}
}

func TestPluginEquals(t *testing.T) {
	var tests = map[string]struct {
		plugin1 *Plugin
		plugin2 *Plugin
		want    bool
	}{
		"Different Endpoints": {
			plugin1: &Plugin{
				Endpoint: "/endpoint1",
			},
			plugin2: &Plugin{
				Endpoint: "/endpoint2",
			},
			want: false,
		},
		"Different Query String": {
			plugin1: &Plugin{
				Endpoint:    "/endpoint1",
				QueryString: "query=test1",
			},
			plugin2: &Plugin{
				Endpoint:    "/endpoint1",
				QueryString: "query=test2",
			},
			want: false,
		},
		"Different Follow Redirects Bool": {
			plugin1: &Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
			},
			plugin2: &Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: false,
			},
			want: false,
		},
		"Equals Checks": {
			plugin1: &Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*Check{
					FakeCheckStatusCode200,
				},
			},
			plugin2: &Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*Check{
					FakeCheckStatusCode200,
				},
			},
			want: true,
		},
		"Not Equals Checks": {
			plugin1: &Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*Check{
					FakeCheckStatusCode200,
					FakeCheckMatchAll,
				},
			},
			plugin2: &Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*Check{
					FakeCheckStatusCode500,
					FakeCheckMatchAll,
				},
			},
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := tc.plugin1.Equals(tc.plugin2)
			if have != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}

func TestSignaturesEquals(t *testing.T) {
	var tests = map[string]struct {
		signatures1 *Signatures
		signatures2 *Signatures
		want        bool
	}{
		"Different Length": {
			signatures1: &Signatures{Plugins: []*Plugin{
				FakePlugin,
				FakeQueryPlugin,
				FakeFollowRedirectPlugin,
			}},
			signatures2: NewSignatures(),
			want:        false,
		},

		"Not the same plugin content": {
			signatures1: &Signatures{Plugins: []*Plugin{
				FakePlugin2,
				FakeQueryPlugin,
				FakeFollowRedirectPlugin,
			}},
			signatures2: &Signatures{Plugins: []*Plugin{
				FakePlugin,
				FakeQueryPlugin,
				FakeFollowRedirectPlugin,
			}},
			want: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := tc.signatures1.Equals(tc.signatures2)
			if have != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}
func TestCheckEquals(t *testing.T) {
	// TODO faire et traiter tous les cas (verifier avec go tool que ca passe dans tous les if)
	var tests = map[string]struct {
		check1 *Check
		check2 *Check
		want   bool
	}{
		"MustMatchOne not Equals": {
			check1: &Check{MustMatchOne: []string{"MATCHONE", "MATCHTWO"}},
			check2: &Check{MustMatchOne: []string{"MATCHFOUR", "MATCHTHREE"}},
			want:   false,
		},
		"MustMatchAll not Equals": {
			check1: &Check{MustMatchAll: []string{"MATCHONE", "MATCHTWO"}},
			check2: &Check{MustMatchAll: []string{"MATCHONE", "MATCHTHREE"}},
			want:   false,
		},
		"MustNotMatch Equals": {
			check1: &Check{MustNotMatch: []string{"MATCHONE", "MATCHTWO"}},
			check2: &Check{MustNotMatch: []string{"MATCHONE", "MATCHTHREE"}},
			want:   false,
		},
		"Name not Equals": {
			check1: &Check{Name: "Name1"},
			check2: &Check{Name: "Name2"},
			want:   false,
		},
		"Remediation not Equals": {
			check1: &Check{Remediation: "ಠ_ಠ"},
			check2: &Check{Remediation: "(°_o)"},
			want:   false,
		},
		"Severity not Equals": {
			check1: &Check{Severity: "High"},
			check2: &Check{Severity: "Medium"},
			want:   false,
		},
		"Description not Equals": {
			check1: &Check{Description: "ಠ_ಠ"},
			check2: &Check{Description: "(°_o)"},
			want:   false,
		},
		"Headers not Equals": {
			check1: &Check{Headers: []string{"Header:OK"}},
			check2: &Check{Headers: []string{"Header:notOK"}},
			want:   false,
		},
		"noHeaders not Equals": {
			check1: &Check{NoHeaders: []string{"NoHeader:OK"}},
			check2: &Check{NoHeaders: []string{"NoHeader:notOK"}},
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := tc.check1.Equals(tc.check2)
			if have != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}

func TestSliceStringEqual(t *testing.T) {
	var tests = map[string]struct {
		slice1 []string
		slice2 []string
		want   bool
	}{
		"Same slices": {
			slice1: []string{"a", "b"},
			slice2: []string{"a", "b"},
			want:   true,
		},
		"Different slices with same length": {
			slice1: []string{"a", "b"},
			slice2: []string{"x", "b"},
			want:   false,
		},
		"Different slices with different length": {
			slice1: []string{"a", "b"},
			slice2: []string{"a", "b", "c"},
			want:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := sliceStringEqual(tc.slice1, tc.slice2)
			if have != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}
