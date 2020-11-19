package core_test

import (
	"context"
	"gochopchop/core"
	"gochopchop/mock"
	"testing"
	"time"
)

// TODO : Test fonctionnel
// TODO : integrer tests unitaires dans la CI (voir github workflows avec go test)

func TestScan(t *testing.T) {
	var tests = map[string]struct {
		ctx    context.Context
		urls   []string
		output []core.Output
	}{
		"no vulnerabilities found":  {ctx: context.Background(), urls: []string{"http://noproblem"}, output: []core.Output{}},
		"all vulnerabilities found": {ctx: context.Background(), urls: []string{"http://problems"}, output: mock.FakeOutput},
		"context is done":           {ctx: context.Background(), urls: []string{"http://noproblem"}, output: []core.Output{}},
		"fetcher problem":           {ctx: context.Background(), urls: []string{"http://unknown"}, output: []core.Output{}},
		"no HTTP Response":          {ctx: context.Background(), urls: []string{"http://nohttpresponse"}, output: []core.Output{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			if name == "context is done" {
				ctx, cancel := context.WithDeadline(tc.ctx, time.Now().Add(-7*time.Hour))
				tc.ctx = ctx
				cancel()
			}

			output, _ := mock.FakeScanner.Scan(tc.ctx, tc.urls)

			for _, haveOutput := range tc.output {
				found := false
				for _, wantOutput := range output {
					if wantOutput.Name == haveOutput.Name {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected: %v, got: %v", tc.output, output)
				}
			}
		})
	}
}
