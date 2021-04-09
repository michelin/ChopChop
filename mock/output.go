package mock

import (
	"github.com/michelin/gochopchop/core"
)

var FakeOutputStatusCode = core.Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckStatusCode200.Name,
	Severity:    FakeCheckStatusCode200.Severity,
	Remediation: FakeCheckStatusCode200.Remediation,
}

var FakeOutputMatchOne = core.Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckMatchOne.Name,
	Severity:    FakeCheckMatchOne.Severity,
	Remediation: FakeCheckMatchOne.Remediation,
}
var FakeOutputMatchAll = core.Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckMatchAll.Name,
	Severity:    FakeCheckMatchAll.Severity,
	Remediation: FakeCheckMatchAll.Remediation,
}

var FakeOutputNotMatch = core.Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckNotMatch.Name,
	Severity:    FakeCheckNotMatch.Severity,
	Remediation: FakeCheckNotMatch.Remediation,
}

var FakeOutputNoHeaders = core.Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckNoHeaders.Name,
	Severity:    FakeCheckNoHeaders.Severity,
	Remediation: FakeCheckNoHeaders.Remediation,
}

var FakeOutputHeaders = core.Output{
	URL:         "http://problems",
	Endpoint:    FakePlugin.Endpoint,
	Name:        FakeCheckHeaders.Name,
	Severity:    FakeCheckHeaders.Severity,
	Remediation: FakeCheckHeaders.Remediation,
}

var FakeOutput = []core.Output{
	FakeOutputStatusCode,
	FakeOutputHeaders,
	FakeOutputNoHeaders,
	FakeOutputMatchAll,
	FakeOutputMatchOne,
	FakeOutputNotMatch,
}

var FakeOutputAsCSV = "url,endpoint,severity,checkName,remediation\nhttp://problems,/,Medium,StatusCode200,uninstall\nhttp://problems,/,High,Headers,uninstall\nhttp://problems,/,Low,NoHeaders,uninstall\nhttp://problems,/,Informational,MustMatchAll,uninstall\nhttp://problems,/,Low,MustMatchOne,uninstall\nhttp://problems,/,High,MustNotMatch,uninstall\n"
var FakeOutputAsTable = "+-----------------+----------+---------------+---------------+-------------+\n| URL             | ENDPOINT | SEVERITY      | PLUGIN        | REMEDIATION |\n+-----------------+----------+---------------+---------------+-------------+\n| http://problems | /        | \x1b[31mHigh\x1b[0m          | Headers       | uninstall   |\n| http://problems | /        | \x1b[31mHigh\x1b[0m          | MustNotMatch  | uninstall   |\n| http://problems | /        | \x1b[32mLow\x1b[0m           | NoHeaders     | uninstall   |\n| http://problems | /        | \x1b[32mLow\x1b[0m           | MustMatchOne  | uninstall   |\n| http://problems | /        | \x1b[33mMedium\x1b[0m        | StatusCode200 | uninstall   |\n| http://problems | /        | \x1b[36mInformational\x1b[0m | MustMatchAll  | uninstall   |\n+-----------------+----------+---------------+---------------+-------------+\n"
var FakeOutputAsJSON = "[{\"url\":\"http://problems\",\"endpoint\":\"/\",\"checkName\":\"StatusCode200\",\"severity\":\"Medium\",\"remediation\":\"uninstall\"},{\"url\":\"http://problems\",\"endpoint\":\"/\",\"checkName\":\"Headers\",\"severity\":\"High\",\"remediation\":\"uninstall\"},{\"url\":\"http://problems\",\"endpoint\":\"/\",\"checkName\":\"NoHeaders\",\"severity\":\"Low\",\"remediation\":\"uninstall\"},{\"url\":\"http://problems\",\"endpoint\":\"/\",\"checkName\":\"MustMatchAll\",\"severity\":\"Informational\",\"remediation\":\"uninstall\"},{\"url\":\"http://problems\",\"endpoint\":\"/\",\"checkName\":\"MustMatchOne\",\"severity\":\"Low\",\"remediation\":\"uninstall\"},{\"url\":\"http://problems\",\"endpoint\":\"/\",\"checkName\":\"MustNotMatch\",\"severity\":\"High\",\"remediation\":\"uninstall\"}]"
