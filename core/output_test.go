package core

var FakeOutputStatusCode = Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckStatusCode.Name,
	Severity:    FakeCheckStatusCode.Severity,
	Remediation: FakeCheckStatusCode.Remediation,
}

var FakeOutputMatchOne = Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckMatchOne.Name,
	Severity:    FakeCheckMatchOne.Severity,
	Remediation: FakeCheckMatchOne.Remediation,
}
var FakeOutputMatchAll = Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckMatchAll.Name,
	Severity:    FakeCheckMatchAll.Severity,
	Remediation: FakeCheckMatchAll.Remediation,
}

var FakeOutputNotMatch = Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckNotMatch.Name,
	Severity:    FakeCheckNotMatch.Severity,
	Remediation: FakeCheckNotMatch.Remediation,
}

var FakeOutputNoHeaders = Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckNoHeaders.Name,
	Severity:    FakeCheckNoHeaders.Severity,
	Remediation: FakeCheckNoHeaders.Remediation,
}

var FakeOutputHeaders = Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckHeaders.Name,
	Severity:    FakeCheckHeaders.Severity,
	Remediation: FakeCheckHeaders.Remediation,
}

var FakeOutput = []Output{
	FakeOutputStatusCode,
	FakeOutputHeaders,
	FakeOutputNoHeaders,
	FakeOutputMatchAll,
	FakeOutputMatchOne,
	FakeOutputNotMatch,
}
