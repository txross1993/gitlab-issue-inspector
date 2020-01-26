package gitlab

import (
	"time"
)

// User represents a user of GitLab
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// TimeStats represents the nested structure of an Issue related to time estimation and time spent toward and issue
type TimeStats struct {
	TimeEstimate   int `json:"time_estimate"`
	TotalTimeSpent int `json:"total_time_spent"`
}

// TaskCompletionStatus represents the nested structure of an Issue related to number of tasks
type TaskCompletionStatus struct {
	Count          int `json:"count"`
	CompletedCount int `json:"completed_count"`
}

// Issue represents the struct for a GitLab issue
type Issue struct {
	ID                   int                  `json:"id"`
	IID                  int                  `json:"iid"`
	Author               User                 `json:"author"`
	Users                []User               `json:"assignees"`
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

// Project represents the struct for a GitLab project
type Project struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Link  string `json:"web_url"`
	Users []int  // Slice of UserIDs; No json binding for users as this is a separate API request
}
