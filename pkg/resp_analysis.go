package pkg

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//ResponseAnalysis of HTTP Request with checks
func ResponseAnalysis(resp *http.Response, statusCode *int32, match []*string, allMatch []*string, noMatch []*string, headers []*string) bool {

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	if statusCode != nil {
		tmpCode := *statusCode
		if int32(resp.StatusCode) != tmpCode {
			return false
		}
	}
	// all element needs to be found
	if allMatch != nil {
		for i := 0; i < len(allMatch); i++ {
			if !strings.Contains(bodyString, *allMatch[i]) {
				return false
			}
		}
	}

	// one elements needs to be found
	if match != nil {
		found := false
		for i := 0; i < len(match); i++ {
			if strings.Contains(bodyString, *match[i]) {
				found = true
			}
		}
		if !found {
			return false
		}
	}

	// if 1 element of list is not found
	if noMatch != nil {
		for i := 0; i < len(noMatch); i++ {
			if strings.Contains(bodyString, *noMatch[i]) {
				return false
			}
		}
	}
	if headers != nil {
		for i := 0; i < len(headers); i++ {
			// Parse headers
			pHeaders := strings.Split(*headers[i], ":")
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
