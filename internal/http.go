package internal

import "net/http"

type HTTPResponse struct {
	StatusCode int
	Body       string
	Header     http.Header
}
