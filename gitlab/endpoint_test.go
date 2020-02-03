package gitlab

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestBuilder(t *testing.T) {

	type queryFunc func() string

	tests := map[string]struct {
		ProjectID int
		IssueID   int
		IssueIID  int
		UserID    int
		Filters   []string
		Func      queryFunc
		Want      string
	}{
		"Issues": {
			Want:    fmt.Sprintf("%s/issues?scope=all&updated_after=2020-01-21T00:22:36.156Z&labels=foo", os.Getenv("BASE_URL")),
			Filters: []string{"scope=all", "updated_after=2020-01-21T00:22:36.156Z", "labels=foo"},
			Func: func() string {
				return GetIssues("2020-01-21T00:22:36.156Z", "foo")
			},
		},
		"Note 262": {
			Want:      fmt.Sprintf("%s/projects/%d/issues/%d/notes?order_by=updated_at&sort=desc", os.Getenv("BASE_URL"), 123, 1),
			ProjectID: 123,
			IssueIID:  1,
			Filters:   []string{"order_by=updated_at", "sort=desc"},
			Func: func() string {
				return GetNotes(123, 1)
			},
		},
		"Project 123": {
			Want:      fmt.Sprintf("%s/projects/%d", os.Getenv("BASE_URL"), 123),
			ProjectID: 123,
			Func: func() string {
				return GetProject(123)
			},
		},
		"ProjectMemebers": {
			Want:      fmt.Sprintf("%s/projects/%d/members", os.Getenv("BASE_URL"), 123),
			ProjectID: 123,
			Func: func() string {
				return GetProjectMembers(123)
			},
		},
		"User 777": {
			Want:   fmt.Sprintf("%s/users/%d", os.Getenv("BASE_URL"), 777),
			UserID: 777,
			Func: func() string {
				return GetUser(777)
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		got := test.Func()
		if test.Want != got {
			t.Errorf("GOT: %s, WANT: %s", got, test.Want)
		}
	}
}

// Sanity check endpoints against known projects
func TestHTTPGet(t *testing.T) {

	if err := godotenv.Load("../.env"); err != nil {
		t.Fatalf("Failed to read ../.env: %v", err)
	}

	if baseUrl := os.Getenv("BASE_URL"); baseUrl == "" {
		t.Fatal("Please set the environment variable BASE_URL to run this test")
	}

	if apiToken := os.Getenv("API_TOKEN"); apiToken == "" {
		t.Fatal("Please set the environment variable API_TOKEN to run this test")
	}

	client := http.Client{}

	type input struct {
		Args []string
		IDs  []int
	}

	tests := map[string]struct {
		Input input
		Func  func(input) string
	}{
		"Issues": {
			Input: input{
				Args: []string{"2020-01-21T00:22:36.156Z", "etech-reporting"},
			},
			Func: func(i input) string {
				return GetIssues(i.Args[0], i.Args[1])
			},
		},
		"Project": {
			Input: input{
				IDs: []int{8006584},
			},
			Func: func(i input) string {
				return GetProject(i.IDs[0])
			},
		},
		"ProjectMembers": {
			Input: input{
				IDs: []int{8006584},
			},
			Func: func(i input) string {
				return GetProjectMembers(i.IDs[0])
			},
		},
		"Note": {
			Input: input{
				IDs: []int{8006584, 2},
			},
			Func: func(i input) string {
				return GetNotes(i.IDs[0], i.IDs[1])
			},
		},
		"User": {
			Input: input{
				IDs: []int{2720359},
			},
			Func: func(i input) string {
				return GetUser(i.IDs[0])
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)
		req, err := http.NewRequest("GET", test.Func(test.Input), nil)
		req.Header.Set("PRIVATE-TOKEN", os.Getenv("API_TOKEN"))
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("GOT response code: %d, WANT: 200", resp.StatusCode)
		}

	}
}
