package issues

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type GitLabUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type TimeStats struct {
	TimeEstimate   int `json:"time_estimate"`
	TotalTimeSpent int `json:"total_time_spent"`
}

type TaskCompletionStatus struct {
	Count          int `json:"count"`
	CompletedCount int `json:"completed_count"`
}

type Issue struct {
	ID                   int                  `json:"iid"`
	Author               GitLabUser           `json:"author"`
	Assignees            []GitLabUser         `json:"assignees"`
	State                string               `json:"state"`
	ProjectID            int                  `json:"project_id"`
	MilestoneID          int                  `json:"milestone"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	ClosedAt             *time.Time           `json:"closed_at"`
	DueDate              string               `json:"due_date"`
	TaskCompletionStatus TaskCompletionStatus `json:"task_completion_status"`
	TimeStats            TimeStats            `jons:"time_stats"`
}

// Fetch returns the list of issues that have been updated on or after the last known updated timestamp
// The issues fetched depend on the labels provided for the issue filter
func Fetch(client *http.Client, labels string, updatedAt string) ([]Issue, error) {
	b, err := getIssues(client, labels)

	if err != nil {
		return nil, err
	}

	var issues []Issue
	err = json.Unmarshal(b, &issues)
	if err != nil {
		return nil, err
	}

	return filterByUpdatedAt(issues, updatedAt)

}

func getIssues(client *http.Client, labels string) ([]byte, error) {
	url := getIssueUrl(labels)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", os.Getenv("API_TOKEN"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func filterByUpdatedAt(issues []Issue, updatedAt string) ([]Issue, error) {
	t, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}

	filteredIssues := []Issue{}
	for _, issue := range issues {
		if !issue.UpdatedAt.Before(t) {
			filteredIssues = append(filteredIssues, issue)
		}
	}

	return filteredIssues, nil
}

func getIssueUrl(labels string) string {
	// Fetch issues older than that
	baseUrl := os.Getenv("BASE_URL") + "/issues"

	labelsFilter := ""
	if labels != "" {
		labelsFilter = fmt.Sprintf("?labels=%s", labels)
	}

	url := fmt.Sprintf("%s%s", baseUrl, labelsFilter)

	return url
}
