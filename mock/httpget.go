package mock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/michelin/gochopchop/internal/httpget"
)

type FakeNetClient map[string]*http.Response

func (f FakeNetClient) Get(url string) (*http.Response, error) {
	// implements IHTTPClient interface
	if res, ok := f[url]; ok {
		return res, nil
	}
	return nil, fmt.Errorf("could not get url : %s", url)
}

var urls = FakeNetClient{
	"url1": &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("foo"))),
	},
}

var FakeFetcher = &httpget.Fetcher{
	Netclient: urls,
}
