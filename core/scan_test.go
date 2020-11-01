package core_test

import (
	"context"
	"fmt"
	"gochopchop/core"
	"gochopchop/internal"
	"net/http"
	"testing"
)

type FakeFetcher map[string]*internal.HTTPResponse

func (f FakeFetcher) Fetch(url string) (*internal.HTTPResponse, error) {
	if res, ok := f[url]; ok {
		return res, nil
	}
	return nil, fmt.Errorf("could not fetch : %s", url)
}

var MyFakeFetcher = FakeFetcher{
	"http://problems/": &internal.HTTPResponse{
		StatusCode: 200,
		Body:       "MATCHONE lorem ipsum MATCHTWO",
		Header: http.Header{
			"Header": []string{"ok"},
		},
	},
	"http://noproblem/": &internal.HTTPResponse{
		StatusCode: 500,
		Body:       "NOTMATCH",
		Header: http.Header{
			"NoHeader": []string{"ok"},
		},
	},
}

func TestScanURL(t *testing.T) {
	var tests = map[string]struct {
		urls   []string
		output []core.Output
	}{
		"noproblem": {urls: []string{"http://noproblem"}, output: []core.Output{}},
		"problems":  {urls: []string{"http://problems"}, output: FakeOutput},
	}

	scanner := core.NewScanner(MyFakeFetcher, MyFakeFetcher, FakeSignatures, 1)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctx := context.Background()
			output, _ := scanner.Scan(ctx, tc.urls)

			for _, haveOutput := range tc.output {
				found := false
				for _, wantOutput := range output {
					if found {
						break
					}
					if wantOutput.Name == haveOutput.Name {
						found = true
					}
				}
				if !found {
					t.Errorf("expected: %v, got: %v", tc.output, output)
				}
			}
		})
	}
}
