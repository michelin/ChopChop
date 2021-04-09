package core_test

import (
	"testing"

	"github.com/michelin/gochopchop/core"
	"github.com/michelin/gochopchop/mock"
)

func TestFilterBySeverity(t *testing.T) {
	var tests = map[string]struct {
		have     *core.Signatures
		want     *core.Signatures
		severity string
	}{
		"Filter nothing": {
			have:     &core.Signatures{Plugins: []*core.Plugin{mock.FakePlugin}},
			want:     &core.Signatures{Plugins: []*core.Plugin{mock.FakePlugin}},
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
		have  *core.Signatures
		want  *core.Signatures
		names []string
	}{
		"Filter nothing": {
			have:  &core.Signatures{Plugins: []*core.Plugin{mock.FakeQueryPlugin}},
			want:  &core.Signatures{Plugins: []*core.Plugin{mock.FakeQueryPlugin}},
			names: []string{mock.FakeCheckStatusCode200.Name},
		},
		"Filter one element": {
			have:  &core.Signatures{Plugins: []*core.Plugin{mock.FakeQueryPlugin}},
			want:  &core.Signatures{},
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
		plugin1 *core.Plugin
		plugin2 *core.Plugin
		want    bool
	}{
		"Different Endpoints": {
			plugin1: &core.Plugin{
				Endpoint: "/endpoint1",
			},
			plugin2: &core.Plugin{
				Endpoint: "/endpoint2",
			},
			want: false,
		},
		"Different Query String": {
			plugin1: &core.Plugin{
				Endpoint:    "/endpoint1",
				QueryString: "query=test1",
			},
			plugin2: &core.Plugin{
				Endpoint:    "/endpoint1",
				QueryString: "query=test2",
			},
			want: false,
		},
		"Different Follow Redirects Bool": {
			plugin1: &core.Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
			},
			plugin2: &core.Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: false,
			},
			want: false,
		},
		"Equals Checks": {
			plugin1: &core.Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*core.Check{
					mock.FakeCheckStatusCode200,
				},
			},
			plugin2: &core.Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*core.Check{
					mock.FakeCheckStatusCode200,
				},
			},
			want: true,
		},
		"Not Equals Checks": {
			plugin1: &core.Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*core.Check{
					mock.FakeCheckStatusCode200,
					mock.FakeCheckMatchAll,
				},
			},
			plugin2: &core.Plugin{
				Endpoint:        "/endpoint1",
				QueryString:     "query=test1",
				FollowRedirects: true,
				Checks: []*core.Check{
					mock.FakeCheckStatusCode500,
					mock.FakeCheckMatchAll,
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
		signatures1 *core.Signatures
		signatures2 *core.Signatures
		want        bool
	}{
		"Different Length": {
			signatures1: &core.Signatures{Plugins: []*core.Plugin{
				mock.FakePlugin,
				mock.FakeQueryPlugin,
				mock.FakeFollowRedirectPlugin,
			}},
			signatures2: core.NewSignatures(),
			want:        false,
		},

		"Not the same plugin content": {
			signatures1: &core.Signatures{Plugins: []*core.Plugin{
				mock.FakePlugin2,
				mock.FakeQueryPlugin,
				mock.FakeFollowRedirectPlugin,
			}},
			signatures2: &core.Signatures{Plugins: []*core.Plugin{
				mock.FakePlugin,
				mock.FakeQueryPlugin,
				mock.FakeFollowRedirectPlugin,
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
	var tests = map[string]struct {
		check1 *core.Check
		check2 *core.Check
		want   bool
	}{
		"MustMatchOne not Equals": {
			check1: &core.Check{MustMatchOne: []string{"MATCHONE", "MATCHTWO"}},
			check2: &core.Check{MustMatchOne: []string{"MATCHFOUR", "MATCHTHREE"}},
			want:   false,
		},
		"MustMatchAll not Equals": {
			check1: &core.Check{MustMatchAll: []string{"MATCHONE", "MATCHTWO"}},
			check2: &core.Check{MustMatchAll: []string{"MATCHONE", "MATCHTHREE"}},
			want:   false,
		},
		"MustNotMatch Equals": {
			check1: &core.Check{MustNotMatch: []string{"MATCHONE", "MATCHTWO"}},
			check2: &core.Check{MustNotMatch: []string{"MATCHONE", "MATCHTHREE"}},
			want:   false,
		},
		"Name not Equals": {
			check1: &core.Check{Name: "Name1"},
			check2: &core.Check{Name: "Name2"},
			want:   false,
		},
		"Remediation not Equals": {
			check1: &core.Check{Remediation: "ಠ_ಠ"},
			check2: &core.Check{Remediation: "(°_o)"},
			want:   false,
		},
		"Severity not Equals": {
			check1: &core.Check{Severity: "High"},
			check2: &core.Check{Severity: "Medium"},
			want:   false,
		},
		"Description not Equals": {
			check1: &core.Check{Description: "ಠ_ಠ"},
			check2: &core.Check{Description: "(°_o)"},
			want:   false,
		},
		"Headers not Equals": {
			check1: &core.Check{Headers: []string{"Header:OK"}},
			check2: &core.Check{Headers: []string{"Header:notOK"}},
			want:   false,
		},
		"NoHeaders not Equals": {
			check1: &core.Check{NoHeaders: []string{"NoHeader:OK"}},
			check2: &core.Check{NoHeaders: []string{"NoHeader:notOK"}},
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
			have := core.SliceStringEqual(tc.slice1, tc.slice2)
			if have != tc.want {
				t.Errorf("expected: %v, got: %v", tc.want, have)
			}
		})
	}
}
