package source

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// getProject returns the project related data from GitLab projects API
func getProject(client *http.Client, projectID int) (data.Project, error) {

	b, err := performRequest(url, client)
	if err != nil {
		return data.Project{}, err
	}

	var project data.Project
	if err := json.Unmarshal(b, &project); err != nil {
		return data.Project{}, err
	}

	ids, err := getProjectUserIDs(endpoint, client)
	if err != nil {
		return data.Project{}, err
	}

	project.UserIDs = ids
	return project, nil
}

func getProjectUserIDs(projectEndpoint string, client *http.Client) ([]int, error) {

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

// getUser returns a user's info from the GitLab API
func getUser(client *http.Client, userID int) (data.User, error) {
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

// getIssues returns the list of issues that have been updated on or after the last known updated timestamp
// The issues geted depend on the labels provided for the issue filter
func getIssues(client *http.Client, updatedAt string, labels string) ([]data.Issue, error) {
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

func getNotes(projectID int, issueIID int, client *http.Client) ([]data.Note, error) {
	url := getNotesUrl(projectID, issueIID)
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
