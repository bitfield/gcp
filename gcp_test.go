package gcp

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"google.golang.org/api/googleapi"
)

var update = flag.Bool("update", false, "Update golden file test fixtures")

func TestInterpretGoogleAPIError(t *testing.T) {
	tests := []struct {
		input error
		want  string
	}{
		{input: error(&googleapi.Error{Code: http.StatusForbidden}), want: "project is not API-enabled"},
		{input: error(&googleapi.Error{Code: http.StatusNotFound}), want: "project not found"},
		{input: error(&googleapi.Error{Code: http.StatusInternalServerError}), want: "API call failed: googleapi: got HTTP response code 500 with body: "},
		{input: error(&googleapi.Error{Code: http.StatusBadRequest}), want: "zone not found"},
		{input: error(errors.New("bogus error")), want: "bogus error"},
	}
	for _, c := range tests {
		got := interpretGoogleAPIError(c.input)
		if got.Error() != c.want {
			t.Errorf("interpretGoogleAPIError(%v) = %v, want %v\n", c.input, got, c.want)
		}
	}
}

func TestJSON2HCL(t *testing.T) {
	tests := []struct {
		inputFile  string
		goldenFile string
	}{
		{inputFile: "instance.json", goldenFile: "instance.hcl"},
		{inputFile: "firewall.json", goldenFile: "firewall.hcl"},
		{inputFile: "network.json", goldenFile: "network.hcl"},
	}
	for _, c := range tests {
		json, err := ioutil.ReadFile(filepath.Join("testdata", c.inputFile))
		if err != nil {
			t.Fatalf("couldn't read test fixture: %s\n", err)
		}
		want_hcl, err := ioutil.ReadFile(filepath.Join("testdata", c.goldenFile))
		if err != nil {
			t.Fatalf("couldn't read test fixture: %s\n", err)
		}
		hcl, err := JSON2HCL(json)
		if err != nil {
			t.Errorf("JSON2HCL(%s) failed: %s\n", c.inputFile, err)
		}
		got_hcl := []byte(hcl)
		if *update {
			fmt.Printf("writing golden file %s: %s\n", c.goldenFile, got_hcl)
			err = ioutil.WriteFile(filepath.Join("testdata", c.goldenFile), []byte(got_hcl), 0644)
			if err != nil {
				t.Fatalf("couldn't update test fixture: %s\n", err)
			}
			want_hcl = got_hcl
		}
		diffText := unifiedDiff(got_hcl, want_hcl)
		if len(diffText) > 0 {
			t.Errorf("JSON2HCL(%s) differs from golden file %s: %s\n", c.inputFile, c.goldenFile, diffText)
		}
	}
}

func unifiedDiff(a, b []byte) string {
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(a)),
		B:       difflib.SplitLines(string(b)),
		Context: 0,
	}
	diffText, _ := difflib.GetUnifiedDiffString(diff)
	return diffText
}
