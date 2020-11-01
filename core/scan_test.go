package core_test

import (
	"context"
	"fmt"
	"gochopchop/core"
	"gochopchop/internal"
	"net/http"
	"reflect"
	"testing"
)

var FakeUrl string = "http://testing"

type FakeFetcher map[string]*internal.HTTPResponse

func (f FakeFetcher) Fetch(url string) (*internal.HTTPResponse, error) {
	if res, ok := f[url]; ok {
		return res, nil
	}
	return nil, fmt.Errorf("could not fetch : %s", url)
}

var MyFakeFetcher = FakeFetcher{
	"http://testing/200": &internal.HTTPResponse{
		StatusCode: 200,
		Body:       "",
		Header:     http.Header{},
	},
	"http://testing/404": &internal.HTTPResponse{
		StatusCode: 404,
		Body:       "",
		Header:     http.Header{},
	},
}

func TestScanURL(t *testing.T) {
	var tests = map[string]struct {
		urls   []string
		output []core.Output
	}{
		"url 200 exists": {urls: []string{FakeUrl}, output: FakeOutput200},
		"url 404 exists":        {urls: []string{FakeUrl}, output: FakeOutput404},
		"url cannot be reached": {urls: []string{FakeUrl}, output: []core.Output{}},
	}

	scanner := core.NewScanner(MyFakeFetcher, MyFakeFetcher, FakeSignatures, 1)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			ctx := context.Background()
			output, _ := scanner.Scan(ctx, tc.urls)

			if !reflect.DeepEqual(output, tc.output) {
				t.Errorf("expected: %v, got: %v", tc.output, output)
			}
		})
	}
}
