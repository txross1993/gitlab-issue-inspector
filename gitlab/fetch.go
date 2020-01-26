package gitlab

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	issuesEndpoint   = "/issues"
	projectsEndpoint = "/projects"
	usersEndpoint    = "/users"
)

// Fetch Project returns the project related data from GitLab projects API
func FetchProject(client *http.Client, projectID string) (Project, error) {
	endpoint := fmt.Sprintf("%s/%s", projectsEndpoint, projectID)
	url := buildQuery(endpoint)
	b, err := performRequest(url, client)
	if err != nil {
		return Project{}, err
	}

	var project Project
	if err := json.Unmarshal(b, &project); err != nil {
		return Project{}, err
	}

	ids, err := fetchProjectUserIDs(endpoint, client)
	if err != nil {
		return Project{}, err
	}

	project.Users = ids
	return project, nil
}

func fetchProjectUserIDs(projectEndpoint string, client *http.Client) ([]int, error) {
	endpoint := fmt.Sprintf("%s/members", projectEndpoint)
	url := buildQuery(endpoint)
	b, err := performRequest(url, client)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(b, &users); err != nil {
		return nil, err
	}

	var userIDs []int
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	return userIDs, nil
}

// FetchUser returns a user's info from the GitLab API
func FetchUser(client *http.Client, userID string) (User, error) {
	endpoint := fmt.Sprintf("%s/%s", usersEndpoint, userID)
	url := buildQuery(endpoint)
	b, err := performRequest(url, client)
	if err != nil {
		return User{}, err
	}

	var user User
	if err := json.Unmarshal(b, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

// FetchIssues returns the list of issues that have been updated on or after the last known updated timestamp
// The issues fetched depend on the labels provided for the issue filter
func FetchIssues(client *http.Client, updatedAt string, labels string) ([]Issue, error) {
	url := getIssuesUrl(updatedAt, labels)
	b, err := performRequest(url, client)

	if err != nil {
		return nil, err
	}

	var issues []Issue
	if err := json.Unmarshal(b, &issues); err != nil {
		return nil, err
	}

	return issues, nil

}

func getIssuesUrl(updatedAt, labels string) string {
	scopeFilter := "scope=all"

	updatedAtFilter := ""
	if updatedAt != "" {
		updatedAtFilter = fmt.Sprintf("updated_after=%s", updatedAt)
	}

	labelsFilter := ""
	if labels != "" {
		labelsFilter = fmt.Sprintf("labels=%s", labels)
	}

	return buildQuery(issuesEndpoint, scopeFilter, updatedAtFilter, labelsFilter)
}
