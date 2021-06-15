package internal_test

import (
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/michelin/gochopchop/internal"
	"gopkg.in/yaml.v2"
)

var (
	statuscode1 int = 1
	statuscode2 int = 2

	triggerStr string = "trigger"
	triggerBts []byte = []byte(triggerStr)
)

func TestCheckMatch(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Check        *internal.Check
		Resp         *internal.HTTPResponse
		ExpectedBool bool
		ExpectedErr  error
	}{
		"nil-check": {
			Check:        nil,
			Resp:         nil,
			ExpectedBool: false,
			ExpectedErr:  &internal.ErrNilParameter{"check"},
		},
		"nil-check-statuscode": {
			Check: &internal.Check{
				StatusCode: nil,
			},
			Resp:         nil,
			ExpectedBool: false,
			ExpectedErr:  &internal.ErrNilParameter{"check.StatusCode"},
		},
		"nil-resp": {
			Check: &internal.Check{
				StatusCode: &statuscode1,
			},
			Resp:         nil,
			ExpectedBool: false,
			ExpectedErr:  &internal.ErrNilParameter{"resp"},
		},
		"different-check-statuscode": {
			Check: &internal.Check{
				StatusCode: &statuscode1,
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode2,
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"must-match-all-false": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchAll: []string{triggerStr},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       []byte{},
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"must-match-one-not-found": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       []byte{},
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"must-match-one-found": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
			},
			ExpectedBool: true,
			ExpectedErr:  nil,
		},
		"must-not-match-false": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				MustNotMatch: []string{triggerStr},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"invalid-header-format": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				Headers:      []string{"Fake-Header:first:second"},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
			},
			ExpectedBool: false,
			ExpectedErr:  &internal.ErrInvalidHeaderFormat{"Fake-Header:first:second"},
		},
		"unknown-header": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				Headers:      []string{"Fake-Header:fake-content"},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
				Header:     http.Header{},
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"valid-match-header": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				Headers:      []string{"Fake-Header:fake-content"},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
				Header: http.Header{
					"Fake-Header": []string{"fake-content"},
				},
			},
			ExpectedBool: true,
			ExpectedErr:  nil,
		},
		"not-match-header": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				Headers:      []string{"Fake-Header:fake-content"},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
				Header: http.Header{
					"Fake-Header": []string{"invalid-content"},
				},
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"invalid-no-header-format": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				NoHeaders:    []string{"Fake-Header:first:second"},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
			},
			ExpectedBool: false,
			ExpectedErr:  &internal.ErrInvalidHeaderFormat{"Fake-Header:first:second"},
		},
		"match-no-header": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
				NoHeaders:    []string{"Fake-Header:fake-content"},
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
				Header: http.Header{
					"Fake-Header": []string{"fake-content"},
				},
			},
			ExpectedBool: false,
			ExpectedErr:  nil,
		},
		"valid-check": {
			Check: &internal.Check{
				StatusCode:   &statuscode1,
				MustMatchOne: []string{triggerStr}, // To pass the MatchOne check
			},
			Resp: &internal.HTTPResponse{
				StatusCode: statuscode1,
				Body:       triggerBts,
			},
			ExpectedBool: true,
			ExpectedErr:  nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			b, err := tt.Check.Match(tt.Resp)

			if b != tt.ExpectedBool {
				t.Errorf("Failed to get expected bool value: got \"%t\" instead of \"%t\".", b, tt.ExpectedBool)
			}

			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

func TestParseSignatures(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Reader       io.Reader
		ExpectedSign *internal.Signatures
		ExpectedErr  error
	}{
		"fail-reader": {
			Reader:       &FailingReadCloser{},
			ExpectedSign: nil,
			ExpectedErr:  errFake,
		},
		"invalid-signatures": {
			Reader:       NewFakeReadCloser("invalid-content"),
			ExpectedSign: nil,
			ExpectedErr:  &yaml.TypeError{},
		},
		"empty-description": {
			Reader: NewFakeReadCloser(`
insecure: false
plugins:
  - endpoints:
      - "example.html"
    checks:
      - name: EXAMPLE
        match:
          - "Example:"
        remediation: Remediation example
        description:
        severity: Informational`),
			ExpectedSign: nil,
			ExpectedErr:  &internal.ErrCheckInvalidField{"EXAMPLE", "description"},
		},
		"empty-remediation": {
			Reader: NewFakeReadCloser(`
insecure: false
plugins:
  - endpoints:
      - "example.html"
    checks:
      - name: EXAMPLE
        match:
          - "Example:"
        remediation:
        description: Remediation example
        severity: Informational`),
			ExpectedSign: nil,
			ExpectedErr:  &internal.ErrCheckInvalidField{"EXAMPLE", "remediation"},
		},
		"empty-severity": {
			Reader: NewFakeReadCloser(`
insecure: false
plugins:
  - endpoints:
      - "example.html"
    checks:
      - name: EXAMPLE
        match:
          - "Example:"
        remediation: Description example
        description: Remediation example
        severity:`),
			ExpectedSign: nil,
			ExpectedErr:  &internal.ErrCheckInvalidField{"EXAMPLE", "severity"},
		},
		"invalid-severity": {
			Reader: NewFakeReadCloser(`
insecure: false
plugins:
  - endpoints:
      - "example.html"
    checks:
      - name: EXAMPLE
        match:
          - "Example:"
        remediation: Description example
        description: Remediation example
        severity: INVALID`),
			ExpectedSign: nil,
			ExpectedErr:  &internal.ErrInvalidSeverity{"INVALID"},
		},
		"invalid-header-format": {
			Reader: NewFakeReadCloser(`
insecure: false
plugins:
  - endpoints:
      - "example.html"
    checks:
      - name: EXAMPLE
        match:
          - "Example:"
        remediation: Description example
        description: Remediation example
        severity: Low
        headers:
          - "Fake-Header:fake:content"`),
			ExpectedSign: nil,
			ExpectedErr:  &internal.ErrInvalidHeaderFormat{Header: "Fake-Header:fake:content"},
		},
		"empty-valid-signatures": {
			Reader:       NewFakeReadCloser(""),
			ExpectedSign: &internal.Signatures{},
			ExpectedErr:  nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			sign, err := internal.ParseSignatures(tt.Reader)

			if !reflect.DeepEqual(sign, tt.ExpectedSign) {
				t.Errorf("Failed to get expected *Signatures: got \"%v\" instead of \"%v\".", sign, tt.ExpectedSign)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

type WriteData interface {
	io.Writer
	Data() []byte
}

type FakeWriter struct {
	data []byte
}

func (f *FakeWriter) Write(data []byte) (int, error) {
	f.data = append(f.data, data...)
	return len(data), nil
}

func (f *FakeWriter) Data() []byte {
	return f.data
}

func NewFakeWriter(data string) *FakeWriter {
	return &FakeWriter{[]byte(data)}
}

var _ = (io.Writer)(&FakeWriter{})
var _ = (WriteData)(&FakeWriter{})

func TestPrintSignatures(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Writer         WriteData
		Sign           *internal.Signatures
		Sev            string
		ExpectedOutput []byte
	}{
		"empty-sign": {
			Writer: NewFakeWriter(""),
			Sign:   &internal.Signatures{},
			Sev:    "",
			ExpectedOutput: []byte(`+----------+------------+--------------+-------------+
| ENDPOINT | CHECK NAME |     SEVERITY | DESCRIPTION |
+----------+------------+--------------+-------------+
+----------+------------+--------------+-------------+
|          |            | TOTAL CHECKS |           0 |
+----------+------------+--------------+-------------+
`),
		},
		"sign-no-matching-severity": {
			Writer: NewFakeWriter(""),
			Sign: &internal.Signatures{
				Plugins: []internal.Plugin{
					{
						Endpoints: []string{"/endpoint1", "/endpoint2"},
						Checks: []internal.Check{
							{
								Name:     "check-1",
								Severity: "known-severity",
							},
						},
						FollowRedirects: false,
					},
				},
			},
			Sev: "unknown-severity",
			ExpectedOutput: []byte(`+----------+------------+--------------+-------------+
| ENDPOINT | CHECK NAME |     SEVERITY | DESCRIPTION |
+----------+------------+--------------+-------------+
+----------+------------+--------------+-------------+
|          |            | TOTAL CHECKS |           0 |
+----------+------------+--------------+-------------+
`),
		},
		"sign-matching-severity": {
			Writer: NewFakeWriter(""),
			Sign: &internal.Signatures{
				Plugins: []internal.Plugin{
					{
						Endpoints: []string{"/endpoint1", "/endpoint2"},
						Checks: []internal.Check{
							{
								Name:     "check-1",
								Severity: "known-severity",
							},
						},
						FollowRedirects: false,
					},
				},
			},
			Sev: "known-severity",
			ExpectedOutput: []byte(`+-------------------------+------------+----------------+-------------+
| ENDPOINT                | CHECK NAME | SEVERITY       | DESCRIPTION |
+-------------------------+------------+----------------+-------------+
| [/endpoint1 /endpoint2] | check-1    | known-severity |             |
+-------------------------+------------+----------------+-------------+
|                         |            | TOTAL CHECKS   | 1           |
+-------------------------+------------+----------------+-------------+
`),
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			internal.PrintSignatures(tt.Sign, tt.Sev, tt.Writer)

			if !reflect.DeepEqual(tt.Writer.Data(), tt.ExpectedOutput) {
				t.Errorf("Failed to get expected output bytes: got \"%v\" instead of \"%v\".", tt.Writer.Data(), tt.ExpectedOutput)
			}
		})
	}
}
