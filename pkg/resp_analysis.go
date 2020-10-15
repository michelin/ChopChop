package pkg

import (
	"bytes"
	"gochopchop/data"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//ResponseAnalysis of HTTP Request with checks
func ResponseAnalysis(resp *http.Response, signature data.Check) bool {

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Restore the io.ReadCloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	bodyString := string(bodyBytes)

	if signature.StatusCode != nil {
		if int32(resp.StatusCode) != *signature.StatusCode {
			return false
		}
	}
	// all element needs to be found
	if signature.AllMatch != nil {
		for i := 0; i < len(signature.AllMatch); i++ {
			if !strings.Contains(bodyString, *signature.AllMatch[i]) {
				return false
			}
		}
	}

	// one elements needs to be found
	if signature.Match != nil {
		found := false
		for i := 0; i < len(signature.Match); i++ {
			if strings.Contains(bodyString, *signature.Match[i]) {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	// if 1 element of list is not found
	if signature.NoMatch != nil {
		for i := 0; i < len(signature.NoMatch); i++ {
			if strings.Contains(bodyString, *signature.NoMatch[i]) {
				return false
			}
		}
	}
	if signature.Headers != nil {
		for i := 0; i < len(signature.Headers); i++ {
			// Parse headers
			pHeaders := strings.Split(*signature.Headers[i], ":")
			if v, kFound := resp.Header[pHeaders[0]]; kFound {
				// Key found - check value
				vFound := false
				for i, n := range v {
					if pHeaders[1] == n {
						_ = i
						vFound = true
					}
				}
				if !vFound {
					return false
				}
			} else {
				return false
			}
		}
	}
	return true
}
