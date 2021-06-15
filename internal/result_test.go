package internal_test

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/michelin/gochopchop/internal"
)

func TestSafeResultsAppend(t *testing.T) {
	t.Parallel()

	s := internal.SafeResults{}
	s.Append(internal.Result{})

	if !cmp.Equal(s.Res, []internal.Result{{}}) {
		t.Error("Failed to properly add a Result in the SafeResults.")
	}
}
func TestSafeResultsGetResults(t *testing.T) {
	t.Parallel()

	resSlice := internal.ResultSlice{
		{
			URL:         "test",
			Endpoint:    "/",
			Name:        "name",
			Severity:    "Informational",
			Remediation: "remediation",
		},
	}

	s := internal.SafeResults{
		Res: resSlice,
	}

	res := s.GetResults()

	if !cmp.Equal(res, resSlice) {
		t.Errorf("Failed to get expected results: \"%v\" instead of \"%v\".", res, resSlice)
	}
}

func TestResultSliceSort(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		ResultSlice         internal.ResultSlice
		ExpectedResultSlice internal.ResultSlice
	}{
		"empty-slice": {
			ResultSlice:         internal.ResultSlice{},
			ExpectedResultSlice: internal.ResultSlice{},
		},
		"not-empty-slice": {
			ResultSlice: internal.ResultSlice{
				{
					URL:      "http://127.0.0.1:8080",
					Endpoint: "/",
				}, {
					URL:      "http://127.0.0.1:8080",
					Endpoint: "/2",
				},
			},
			ExpectedResultSlice: internal.ResultSlice{
				{
					URL:      "http://127.0.0.1:8080",
					Endpoint: "/",
				}, {
					URL:      "http://127.0.0.1:8080",
					Endpoint: "/2",
				},
			},
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			sort.Stable(tt.ResultSlice)

			if !cmp.Equal(tt.ResultSlice, tt.ExpectedResultSlice) {
				t.Errorf("Failed to sort as expected: got \"%v\" instead of \"%v\".", tt.ResultSlice, tt.ExpectedResultSlice)
			}
		})
	}
}
