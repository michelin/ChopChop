package export

import (
	"testing"

	"github.com/michelin/gochopchop/core"
	"github.com/michelin/gochopchop/mock"

	"github.com/spf13/afero"
)

func TestExportCSV(t *testing.T) {
	appfs := afero.Afero{Fs: afero.NewMemMapFs()}
	filename := "formatcsv"

	var tests = map[string]struct {
		output []core.Output
		want   string
	}{
		"correct formatting": {output: mock.FakeOutput, want: mock.FakeOutputAsCSV},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, _ := appfs.Create(filename)
			_ = exportCSV(f, tc.output)
			contents, _ := appfs.ReadFile(filename)
			got := string(contents)
			if got != tc.want {
				t.Errorf("want : %q, got : %q", tc.want, got)
			}
		})
	}
}
func TestExportJSON(t *testing.T) {
	appfs := afero.Afero{Fs: afero.NewMemMapFs()}
	filename := "formatjson"

	var tests = map[string]struct {
		output []core.Output
		want   string
	}{
		"correct formatting": {output: mock.FakeOutput, want: mock.FakeOutputAsJSON},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f, _ := appfs.Create(filename)
			_ = exportJSON(f, tc.output)
			contents, _ := appfs.ReadFile(filename)
			got := string(contents)
			if got != tc.want {
				t.Errorf("want : %q, got : %q", tc.want, got)
			}
		})
	}
}
