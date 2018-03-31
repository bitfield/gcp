package gcp

import (
	"errors"
	"net/http"
	"testing"

	"google.golang.org/api/googleapi"
)

func TestInterpretGoogleAPIError(t *testing.T) {
	tests := []struct {
		input error
		want  string
	}{
		{input: error(&googleapi.Error{Code: http.StatusForbidden}), want: "project is not API-enabled"},
		{input: error(&googleapi.Error{Code: http.StatusNotFound}), want: "project not found"},
		{input: error(&googleapi.Error{Code: http.StatusInternalServerError}), want: "API call failed: googleapi: got HTTP response code 500 with body: "},
		{input: error(errors.New("bogus error")), want: "bogus error"},
	}
	for _, c := range tests {
		got := interpretGoogleAPIError(c.input)
		if got.Error() != c.want {
			t.Errorf("interpretGoogleAPIError(%v) = %v, want %v\n", c.input, got, c.want)
		}
	}
}
