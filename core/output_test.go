package core_test

import "gochopchop/core"

var FakeOutput200 = []core.Output{
	{
		URL:         "http://testing",
		Endpoint:    "/200",
		Name:        "url200",
		Severity:    "Medium",
		Remediation: "uninstall",
	},
}

var FakeOutput404 = []core.Output{
	{
		URL:         "http://testing",
		Endpoint:    "/404",
		Name:        "url404",
		Severity:    "Medium",
		Remediation: "uninstall",
	},
}
