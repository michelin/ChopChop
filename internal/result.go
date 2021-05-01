package internal

import (
	"sort"
	"sync"
)

// Result wraps result for each finding of the scan.
type Result struct {
	URL         string `json:"url"`
	Endpoint    string `json:"endpoint"`
	Name        string `json:"checkName"`
	Severity    string `json:"severity"`
	Remediation string `json:"remediation"`
}

// SafeResults stores a Result slice.
type SafeResults struct {
	mux sync.Mutex
	Res []Result
}

// Append adds a Result to the existing ones. Does not
// check for duplications.
func (s *SafeResults) Append(res Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.Res = append(s.Res, res)
}

// GetResults returns the Result array of SafeResults.
func (s *SafeResults) GetResults() ResultSlice {
	return s.Res
}

// ResultSlice wraps a Result slice to enable sorting it.
type ResultSlice []Result

func (rs ResultSlice) Len() int {
	return len(rs)
}

func (rs ResultSlice) Less(i, j int) bool {
	if rs[i].URL == rs[j].URL {
		return rs[i].Endpoint < rs[j].Endpoint
	}
	return rs[i].URL < rs[j].URL
}

func (rs ResultSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

var _ = (sort.Interface)(ResultSlice{})
