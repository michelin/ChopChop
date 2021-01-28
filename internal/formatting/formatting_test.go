package formatting_test

import (
	"bytes"
	"gochopchop/internal/formatting"
	"gochopchop/mock"
	"testing"
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
