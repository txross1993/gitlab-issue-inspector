package issues

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	td "github.com/txross1993/gitlab-issue-inspector/testdata"
)

// RoundTripper testing functionality: http://hassansin.github.io/Unit-Testing-http-client-in-Go
type RoundTripFunc func(req *http.Request) *http.Response

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req), nil
}

func newTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestFetch(t *testing.T) {
	b := td.HelperReadTestData(t, "issues.json")
	testLabel := "myLabel"
	os.Setenv("BASE_URL", "https://gitlab.com/api/v4")
	client := newTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), fmt.Sprintf("%s/issues%s", os.Getenv("BASE_URL"), fmt.Sprintf("?labels=%s", testLabel)))

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})
	updatedAt := "2020-01-21T00:22:36.156Z"

	got, err := Fetch(client, testLabel, updatedAt)
	if err != nil {
		t.Error(err)
	}

	testIssues := []Issue{}
	if err := json.Unmarshal(b, &testIssues); err != nil {
		t.Error(err)
	}
	expected := []Issue{
		testIssues[0],
	}

	assert.Equal(t, expected, got)
}

func TestGetIssuesUrl(t *testing.T) {
	baseUrl := "https://gitlab.example.com/api/v4"
	os.Setenv("BASE_URL", baseUrl)

	tests := map[string]struct {
		labels string
		want   string
	}{
		"all-inputs-supplied": {
			labels: "etech-reporting,doing",
			want:   fmt.Sprintf("%s/issues?labels=etech-reporting,doing", baseUrl),
		},
		"null-values": {
			want: fmt.Sprintf("%s/issues", baseUrl),
		},
	}

	for name, test := range tests {
		t.Logf("Runing test case: %s", name)
		got := getIssueUrl(test.labels)
		if got != test.want {
			t.Fatalf("GOT %s, WANT %s", got, test.want)
		}
	}

}
