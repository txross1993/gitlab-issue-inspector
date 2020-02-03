package gitlab

import (
	"fmt"
	"os"
)

const (
	issuesEndpoint   = "issues"
	projectsEndpoint = "projects"
	usersEndpoint    = "users"
	notesEndpoint    = "notes"
)

func build(endpoint string, filters ...string) string {
	baseUrl := os.Getenv("BASE_URL") + "/" + endpoint

	for i, f := range filters {
		if f == "" {
			continue
		}
		if i == 0 {
			baseUrl += fmt.Sprintf("?%s", f)
		} else {
			baseUrl += fmt.Sprintf("&%s", f)
		}
	}

	return baseUrl
}
