package main

import (
	"bytes"
	"testing"

	compute "google.golang.org/api/compute/v1"
)

func TestDumpInstances(t *testing.T) {
	instances := []*compute.Instance{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}
	var got, want bytes.Buffer
	want.WriteString(`"name" = "foo""name" = "bar""name" = "baz"`)
	if err := dumpInstances(&got, instances); err != nil {
		t.Errorf("dumpInstances errored: %v", err)
	}
	if !bytes.Equal(want.Bytes(), got.Bytes()) {
		t.Errorf("dumpInstances returned %s, want %s", got.String(), want.String())
	}
}
