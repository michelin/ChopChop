package core_test

import (
	"gochopchop/core"
	"testing"
)

func createInt32(x int32) *int32 {
	return &x
}

// Checks

var fakeCheckStatusCode = &core.Check{
	Name:        "StatusCode",
	Severity:    "Medium",
	Remediation: "uninstall",
	StatusCode:  createInt32(200),
}

var fakeCheckNoHeaders = &core.Check{
	Name:        "NoHeaders",
	Severity:    "Medium",
	Remediation: "uninstall",
	NoHeaders:   []string{"NoHeader:ok"},
}

var fakeCheckHeaders = &core.Check{
	Name:        "Headers",
	Severity:    "Medium",
	Remediation: "uninstall",
	Headers:     []string{"Header:ok"},
}

var fakeCheckMatchOne = &core.Check{
	Name:         "MustMatchOne",
	Severity:     "Medium",
	Remediation:  "uninstall",
	MustMatchOne: []string{"MATCHONE", "MATCHTWO"},
}

var fakeCheckMatchAll = &core.Check{
	Name:         "MustMatchAll",
	Severity:     "Medium",
	Remediation:  "uninstall",
	MustMatchAll: []string{"MATCHONE", "MATCHTWO"},
}

var fakeCheckNotMatch = &core.Check{
	Name:         "MustNotMatch",
	Severity:     "Medium",
	Remediation:  "uninstall",
	MustNotMatch: []string{"NOTMATCH"},
}

// Plugins

var fakePlugin = &core.Plugin{
	Endpoint: "/",
	Checks: []*core.Check{
		fakeCheckStatusCode,
		fakeCheckHeaders,
		fakeCheckNoHeaders,
		fakeCheckMatchAll,
		fakeCheckMatchOne,
		fakeCheckNotMatch,
	},
}

// Signatures

var FakeSignatures = &core.Signatures{
	Plugins: []*core.Plugin{
		fakePlugin,
	},
}

func TestFilterBySeverity(t *testing.T) {
	want := &core.Signatures{}
	have := FakeSignatures
	have.FilterBySeverity("High")
	if !want.Equals(have) {
		t.Errorf("expected: %v, got: %v", want, have)
	}
}

func TestFilterByNames(t *testing.T) {
	want := &core.Signatures{}
	have := FakeSignatures
	have.FilterByNames([]string{"UnknownCheck"})
	if !want.Equals(have) {
		t.Errorf("expected: %v, got: %v", want, have)
	}
}
