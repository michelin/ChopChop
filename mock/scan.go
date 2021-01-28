package mock

import (
	"fmt"
	"gochopchop/core"
	"gochopchop/internal"
	"net/http"
)

var FakeScanner = core.NewScanner(MyFakeFetcher, MyFakeFetcher, FakeSignatures, 1)

type FakeFetcherWithoutNetclient map[string]*internal.HTTPResponse

func (f FakeFetcherWithoutNetclient) Fetch(url string) (*internal.HTTPResponse, error) {
	if res, ok := f[url]; ok {
		return res, nil
	}
	return nil, fmt.Errorf("could not fetch : %s", url)
}

var MyFakeFetcher = FakeFetcherWithoutNetclient{
	"http://problems/": &internal.HTTPResponse{
		StatusCode: 200,
		Body:       "MATCHONE lorem ipsum MATCHTWO",
		Header: http.Header{
			"Header":  []string{"ok"},
			"Header2": []string{"ok"},
		},
	},
	"http://noproblem/": &internal.HTTPResponse{
		StatusCode: 500,
		Body:       "NOTMATCH",
		Header: http.Header{
			"Header":    []string{"pasdutout"},
			"NoHeader":  []string{"ok"},
			"NoHeader2": []string{"ok"},
		},
	},
	"http://noproblem/?query=test": &internal.HTTPResponse{
		StatusCode: 500,
	},
}
