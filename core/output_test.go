package core_test

import "gochopchop/core"

var FakeOutputStatusCode = core.Output{
	URL:         "http://problems",
	Endpoint:    fakePlugin.Endpoint,
	Name:        fakeCheckStatusCode.Name,
	Severity:    fakeCheckStatusCode.Severity,
	Remediation: fakeCheckStatusCode.Remediation,
}

var FakeOutputMatchOne = core.Output{
	URL:         "http://problems",
	Endpoint:    fakePlugin.Endpoint,
	Name:        fakeCheckMatchOne.Name,
	Severity:    fakeCheckMatchOne.Severity,
	Remediation: fakeCheckMatchOne.Remediation,
}
var FakeOutputMatchAll = core.Output{
	URL:         "http://problems",
	Endpoint:    fakePlugin.Endpoint,
	Name:        fakeCheckMatchAll.Name,
	Severity:    fakeCheckMatchAll.Severity,
	Remediation: fakeCheckMatchAll.Remediation,
}

var FakeOutputNotMatch = core.Output{
	URL:         "http://problems",
	Endpoint:    fakePlugin.Endpoint,
	Name:        fakeCheckNotMatch.Name,
	Severity:    fakeCheckNotMatch.Severity,
	Remediation: fakeCheckNotMatch.Remediation,
}

var FakeOutputNoHeaders = core.Output{
	URL:         "http://problems",
	Endpoint:    fakePlugin.Endpoint,
	Name:        fakeCheckNoHeaders.Name,
	Severity:    fakeCheckNoHeaders.Severity,
	Remediation: fakeCheckNoHeaders.Remediation,
}

var FakeOutputHeaders = core.Output{
	URL:         "http://problems",
	Endpoint:    fakePlugin.Endpoint,
	Name:        fakeCheckHeaders.Name,
	Severity:    fakeCheckHeaders.Severity,
	Remediation: fakeCheckHeaders.Remediation,
}

var FakeOutput = []core.Output{
	FakeOutputStatusCode,
	FakeOutputHeaders,
	FakeOutputNoHeaders,
	FakeOutputMatchAll,
	FakeOutputMatchOne,
	FakeOutputNotMatch,
}
