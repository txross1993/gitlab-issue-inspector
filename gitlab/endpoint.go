package gitlab

import (
	"fmt"
)

// GetIssues returns the http endpoint uri string for getting all issues updated after @updatedAt and containing all comma-separated lables in @labels string
func GetIssues(updatedAt string, labels string) string {
	scopeFilter := "scope=all"

	updatedAtFilter := ""
	if updatedAt != "" {
		updatedAtFilter = fmt.Sprintf("updated_after=%s", updatedAt)
	}

	labelsFilter := ""
	if labels != "" {
		labelsFilter = fmt.Sprintf("labels=%s", labels)
	}

	return build(issuesEndpoint, scopeFilter, updatedAtFilter, labelsFilter)
}

// GetNotes returns the http endpoint uri string for getting all notes belonging to @projectID, @issueIID
func GetNotes(projectID, issueIID int) string {
	endpoint := fmt.Sprintf("%s/%d/%s/%d/%s", projectsEndpoint, projectID, issuesEndpoint, issueIID, notesEndpoint)

	orderByUpdatedAt := "order_by=updated_at"
	sortDesc := "sort=desc"

	return build(endpoint, orderByUpdatedAt, sortDesc)
}

// GetProject returns the http endpoint uri string for getting all @projectID information
func GetProject(projectID int) string {
	endpoint := fmt.Sprintf("%s/%d", projectsEndpoint, projectID)
	return build(endpoint)
}

// GetProjectMembers returns the http endpoint uri string for getting all @projectID members list
func GetProjectMembers(projectID int) string {
	endpoint := fmt.Sprintf("%s/%d/members", projectsEndpoint, projectID)
	return build(endpoint)
}

// GetUser returns the http endpoint uri string for getting all @userID information
func GetUser(userID int) string {
	endpoint := fmt.Sprintf("%s/%d", usersEndpoint, userID)
	return build(endpoint)
}
