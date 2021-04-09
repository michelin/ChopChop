package formatting_test

import (
	"bytes"
	"testing"

	"github.com/michelin/gochopchop/internal/formatting"
	"github.com/michelin/gochopchop/mock"
)

func TestFormatOutputTable(t *testing.T) {
	mirror := new(bytes.Buffer)
	output := mock.FakeOutput
	formatting.PrintTable(output, mirror)
	got := mirror.String()
	want := mock.FakeOutputAsTable
	if got != want {
		t.Errorf("want : %q, got : %q", want, got)
	}
}
