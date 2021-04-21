package internal_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/michelin/gochopchop/internal"
)

func TestExportersList(t *testing.T) {
	t.Parallel()

	expL := internal.ExportersList()

	if expL == "" {
		t.Error("Exporters can't be empty")
	}
}

func TestCheckSeverities(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Results     []*internal.Result
		MaxSev      internal.Severity
		ExpectedErr error
	}{
		"invalid-max-severity": {
			Results:     nil,
			MaxSev:      -1,
			ExpectedErr: &internal.ErrUnsupportedSeverity{-1},
		},
		"invalid-severity": {
			Results: []*internal.Result{
				{Severity: "invalid-severity"},
			},
			MaxSev:      internal.Informational,
			ExpectedErr: &internal.ErrInvalidSeverity{"invalid-severity"},
		},
		"reached-severity": {
			Results: []*internal.Result{
				{Severity: "High"},
			},
			MaxSev:      internal.Informational,
			ExpectedErr: &internal.ErrMaxSeverityReached{Max: internal.Informational, Sev: internal.High},
		},
		"valid-severities": {
			Results:     nil,
			MaxSev:      0,
			ExpectedErr: nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			err := internal.CheckSeverities(tt.Results, tt.MaxSev)

			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

type WriteDataCloser interface {
	io.WriteCloser
	Data() []byte
}

type FakeWriteCloser struct {
	data []byte
}

func (f *FakeWriteCloser) Write(data []byte) (int, error) {
	f.data = append(f.data, data...)
	return len(data), nil
}

func (f *FakeWriteCloser) Close() error {
	return nil
}

func (f *FakeWriteCloser) Data() []byte {
	return f.data
}

func NewFakeWriteCloser(data string) *FakeWriteCloser {
	return &FakeWriteCloser{[]byte(data)}
}

type FailingWriteCloser struct {
	io.WriteCloser
}

func (f *FailingWriteCloser) Write(data []byte) (int, error) {
	return 0, errFake
}

func (f *FailingWriteCloser) Close() error {
	return nil
}

func (f *FailingWriteCloser) Data() []byte {
	return nil
}

// TODO improve to test for a Marshal error
func TestExportJSON(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Results        []*internal.Result
		Writer         WriteDataCloser
		ExpectedOutput []byte
		ExpectedErr    error
	}{
		"nil-results": {
			Results:        nil,
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: []byte("null"), // null because Results == nil
			ExpectedErr:    nil,
		},
		"empty-results": {
			Results:        []*internal.Result{},
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: []byte("[]"),
			ExpectedErr:    nil,
		},
		"failing-writer": {
			Results:        nil,
			Writer:         &FailingWriteCloser{},
			ExpectedOutput: nil,
			ExpectedErr:    errFake,
		},
		"valid": {
			Results: []*internal.Result{
				{
					URL:         "https://www.michelin.com/",
					Endpoint:    "/",
					Name:        "EXAMPLE",
					Severity:    "Low",
					Remediation: "Remediate",
				},
			},
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: []byte(`[{"url":"https://www.michelin.com/","endpoint":"/","checkName":"EXAMPLE","severity":"Low","remediation":"Remediate"}]`),
			ExpectedErr:    nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			err := internal.ExportJSON(tt.Results, tt.Writer)

			if !reflect.DeepEqual(tt.Writer.Data(), tt.ExpectedOutput) {
				t.Errorf("Failed to get expected output bytes: got \"%v\" instead of \"%v\".", tt.Writer.Data(), tt.ExpectedOutput)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

func TestExportCSV(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Results        []*internal.Result
		Writer         WriteDataCloser
		ExpectedOutput []byte
		ExpectedErr    error
	}{
		"nil-results": {
			Results:        nil,
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: []byte("url,endpoint,severity,checkName,remediation\n"),
			ExpectedErr:    nil,
		},
		"empty-results": {
			Results:        []*internal.Result{},
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: []byte("url,endpoint,severity,checkName,remediation\n"),
			ExpectedErr:    nil,
		},
		"failing-writer": {
			Results:        nil,
			Writer:         &FailingWriteCloser{},
			ExpectedOutput: nil,
			ExpectedErr:    errFake,
		},
		"valid": {
			Results: []*internal.Result{
				{
					URL:         "https://www.michelin.com/",
					Endpoint:    "/",
					Name:        "EXAMPLE",
					Severity:    "Low",
					Remediation: "Remediate",
				},
			},
			Writer: &FakeWriteCloser{},
			ExpectedOutput: []byte(`url,endpoint,severity,checkName,remediation
https://www.michelin.com/,/,Low,EXAMPLE,Remediate
`), // Notice this \n
			ExpectedErr: nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			err := internal.ExportCSV(tt.Results, tt.Writer)

			if !reflect.DeepEqual(tt.Writer.Data(), tt.ExpectedOutput) {
				t.Errorf("Failed to get expected output bytes: got \"%v\" instead of \"%v\".", tt.Writer.Data(), tt.ExpectedOutput)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}

func TestExportTable(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Results        []*internal.Result
		Writer         WriteDataCloser
		ExpectedOutput []byte
		ExpectedErr    error
	}{
		"invalid-severity": {
			Results: []*internal.Result{
				{
					Severity: "invalid-severity",
				},
			},
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: nil,
			ExpectedErr:    &internal.ErrInvalidSeverity{"invalid-severity"},
		},
		"full-table": {
			Results: []*internal.Result{
				{Severity: "High"},
				{Severity: "Medium"},
				{Severity: "Low"},
				{Severity: "Informational"},
			},
			Writer:         &FakeWriteCloser{},
			ExpectedOutput: []byte("+-----+----------+---------------+--------+-------------+\n| URL | ENDPOINT | SEVERITY      | PLUGIN | REMEDIATION |\n+-----+----------+---------------+--------+-------------+\n|     |          | \033[31mHigh\033[0m          |        |             |\n|     |          | \033[32mLow\033[0m           |        |             |\n|     |          | \033[33mMedium\033[0m        |        |             |\n|     |          | \033[36mInformational\033[0m |        |             |\n+-----+----------+---------------+--------+-------------+\n"),
			ExpectedErr:    nil,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			err := internal.ExportTable(tt.Results, tt.Writer)

			if !reflect.DeepEqual(tt.Writer.Data(), tt.ExpectedOutput) {
				t.Errorf("Failed to get expected output bytes: got \"%s\" instead of \"%s\".", tt.Writer.Data(), tt.ExpectedOutput)
			}
			checkErr(err, tt.ExpectedErr, t)
		})
	}
}
