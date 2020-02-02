package source

import (
	"encoding/json"
	"fmt"
	"net/http"

	data "github.com/txross1993/gitlab-issue-inspector/data"
)

const (
	issuesEndpoint   = "/issues"
	projectsEndpoint = "/projects"
	usersEndpoint    = "/users"
	notesEndpoint    = "/notes"
)

// fetchProject returns the project related data from GitLab projects API
func fetchProject(client *http.Client, projectID int) (data.Project, error) {
	endpoint := fmt.Sprintf("%s/%d", projectsEndpoint, projectID)
	url := buildQuery(endpoint)
	b, err := performRequest(url, client)
	if err != nil {
		return data.Project{}, err
	}

	var project data.Project
	if err := json.Unmarshal(b, &project); err != nil {
		return data.Project{}, err
	}

	ids, err := fetchProjectUserIDs(endpoint, client)
	if err != nil {
		return data.Project{}, err
	}

	project.UserIDs = ids
	return project, nil
}

func fetchProjectUserIDs(projectEndpoint string, client *http.Client) ([]int, error) {
	endpoint := fmt.Sprintf("%s/members", projectEndpoint)
	url := buildQuery(endpoint)
	b, err := performRequest(url, client)
	if err != nil {
		return nil, err
	}

	var users []data.User
	if err := json.Unmarshal(b, &users); err != nil {
		return nil, err
	}

	var userIDs []int
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	return userIDs, nil
}

// fetchUser returns a user's info from the GitLab API
func fetchUser(client *http.Client, userID int) (data.User, error) {
	endpoint := fmt.Sprintf("%s/%d", usersEndpoint, userID)
	url := buildQuery(endpoint)
	b, err := performRequest(url, client)
	if err != nil {
		return data.User{}, err
	}

	var user data.User
	if err := json.Unmarshal(b, &user); err != nil {
		return data.User{}, err
	}

	return user, nil
}

// fetchIssues returns the list of issues that have been updated on or after the last known updated timestamp
// The issues fetched depend on the labels provided for the issue filter
func fetchIssues(client *http.Client, updatedAt string, labels string) ([]data.Issue, error) {
	url := getIssuesUrl(updatedAt, labels)
	b, err := performRequest(url, client)

	if err != nil {
		return nil, err
	}

	var issues []data.Issue
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

func fetchNotes(projectID int, issueIID int, client *http.Client) ([]data.Note, error) {
	url := fetchNotesUrl(projectID, issueIID)
	b, err := performRequest(url, client)
	if err != nil {
		return nil, err
	}
	var notes []data.Note
	if err := json.Unmarshal(b, &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func fetchNotesUrl(projectID int, issueIID int) string {
	endpoint := fmt.Sprintf("%s/%d%s/%d%s", projectsEndpoint, projectID, issuesEndpoint, issueIID, notesEndpoint)

	orderByUpdatedAt := "order_by=updated_at"
	sortDesc := "sort=desc"

	return buildQuery(endpoint, orderByUpdatedAt, sortDesc)
}
