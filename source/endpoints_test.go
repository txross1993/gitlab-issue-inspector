package source

import (
	"fmt"
	"net/http"
	"testing"
)

// RoundTripper testing functionality: http://hassansin.github.io/Unit-Testing-http-client-in-Go
type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestGetIssuesUrl(t *testing.T) {

	tests := map[string]struct {
		updatedAt string
		labels    string
		want      string
	}{
		"all-inputs-supplied": {
			updatedAt: "2020-01-21T00:22:36.156Z",
			labels:    "etech-reporting,doing",
			want:      fmt.Sprintf("%s/issues?scope=all&updated_after=2020-01-21T00:22:36.156Z&labels=etech-reporting,doing", baseUrl),
		},
		"null-values": {
			want: fmt.Sprintf("%s/issues?scope=all", baseUrl),
		},
	}

	for name, test := range tests {
		t.Logf("Runing test case: %s", name)
		got := getIssuesUrl(test.updatedAt, test.labels)
		if got != test.want {
			t.Fatalf("GOT %s, WANT %s", got, test.want)
		}
	}

}
