package mock

import (
	"gochopchop/core"
)

func createInt32(x int32) *int32 {
	return &x
}

// Checks

var FakeCheckStatusCode200 = &core.Check{
	Name:        "StatusCode200",
	Severity:    "Medium",
	Remediation: "uninstall",
	StatusCode:  createInt32(200),
}

var FakeCheckStatusCode500 = &core.Check{
	Name:        "StatusCode500",
	Severity:    "High",
	Remediation: "uninstall",
	StatusCode:  createInt32(500),
}

var FakeCheckNoHeaders = &core.Check{
	Name:        "NoHeaders",
	Severity:    "Low",
	Remediation: "uninstall",
	NoHeaders:   []string{"NoHeader:ok"},
}

var FakeCheckNoHeadersKeyOnly = &core.Check{
	Name:        "NoHeaders",
	Severity:    "Informational",
	Remediation: "uninstall",
	NoHeaders:   []string{"NoHeader2"},
}

var FakeCheckHeaders = &core.Check{
	Name:        "Headers",
	Severity:    "High",
	Remediation: "uninstall",
	Headers:     []string{"Header:ok"},
}
var FakeCheckHeaders2 = &core.Check{
	Name:        "Headers",
	Severity:    "Medium",
	Remediation: "uninstall",
	Headers:     []string{"Header2:ok"},
}

var FakeCheckMatchOne = &core.Check{
	Name:         "MustMatchOne",
	Severity:     "Low",
	Remediation:  "uninstall",
	MustMatchOne: []string{"MATCHONE", "MATCHTWO"},
}

var FakeCheckMatchAll = &core.Check{
	Name:         "MustMatchAll",
	Severity:     "Informational",
	Remediation:  "uninstall",
	MustMatchAll: []string{"MATCHONE", "MATCHTWO"},
}

var FakeCheckNotMatch = &core.Check{
	Name:         "MustNotMatch",
	Severity:     "High",
	Remediation:  "uninstall",
	MustNotMatch: []string{"NOTMATCH"},
}

// Plugins

var FakePlugin = &core.Plugin{
	Endpoint: "/",
	Checks: []*core.Check{
		FakeCheckStatusCode200,
		FakeCheckHeaders,
		FakeCheckHeaders2,
		FakeCheckNoHeaders,
		FakeCheckNoHeadersKeyOnly,
		FakeCheckMatchAll,
		FakeCheckMatchOne,
		FakeCheckNotMatch,
	},
}

var FakeQueryPlugin = &core.Plugin{
	Endpoint:    "/",
	QueryString: "query=test",
	Checks: []*core.Check{
		FakeCheckStatusCode200,
	},
}

var FakePlugin2 = &core.Plugin{
	Endpoint:    "/fake",
	QueryString: "query=test",
	Checks: []*core.Check{
		FakeCheckStatusCode500,
	},
}

var FakeFollowRedirectPlugin = &core.Plugin{
	Endpoint: "/",
	Checks: []*core.Check{
		FakeCheckStatusCode200,
	},
	FollowRedirects: true,
}

// Signatures
var FakeSignatures = &core.Signatures{
	Plugins: []*core.Plugin{
		FakePlugin,
		FakeQueryPlugin,
		FakeFollowRedirectPlugin,
	},
}
