package core

import (
	"testing"
)

func createInt32(x int32) *int32 {
	return &x
}

// Checks

var FakeCheckStatusCode = &Check{
	Name:        "StatusCode",
	Severity:    "Medium",
	Remediation: "uninstall",
	StatusCode:  createInt32(200),
}
var FakeCheckStatusCode2 = &Check{
	Name:        "StatusCode",
	Severity:    "Medium",
	Remediation: "uninstall",
	StatusCode:  createInt32(200),
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
		FakeCheckStatusCode,
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
		FakeCheckStatusCode,
	},
}

var FakeFollowRedirectPlugin = &Plugin{
	Endpoint: "/",
	Checks: []*Check{
		FakeCheckStatusCode2,
	},
	FollowRedirects: true,
}

// Signatures
//TODO FIXME signatures ne prends pas plus d'un seul plugin sinon plante
var FakeSignatures = &Signatures{
	Plugins: []*Plugin{
		FakePlugin,
		//FakeQueryPlugin,
		FakeFollowRedirectPlugin,
	},
}

func TestFilterBySeverity(t *testing.T) {
	want := NewSignatures()
	have := FakeSignatures
	have.FilterBySeverity("High")
	if !want.Equals(have) {
		t.Errorf("expected: %v, got: %v", want, have)
	}
}

func TestFilterByNames(t *testing.T) {
	// TODO faire et traiter tous les cas
	want := NewSignatures()
	have := FakeSignatures
	have.FilterByNames([]string{"UnknownCheck"})
	if !want.Equals(have) {
		t.Errorf("expected: %v, got: %v", want, have)
	}
}

func TestPluginEquals(t *testing.T) {
	// TODO faire et traiter tous les cas
	// TODO faire et traiter tous les cas (verifier avec go tool que ca passe dans tous les if)
	var tests = map[string]struct {
		plugin1 *Plugin
		plugin2 *Plugin
		want    bool
	}{
		"MustMatchOne": {
			plugin1: &Plugin{},
			plugin2: &Plugin{},
			want:    true,
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

func TestCheckEquals(t *testing.T) {
	// TODO faire et traiter tous les cas (verifier avec go tool que ca passe dans tous les if)
	var tests = map[string]struct {
		check1 *Check
		check2 *Check
		want   bool
	}{
		"MustMatchOne": {
			check1: &Check{},
			check2: &Check{},
			want:   true,
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
