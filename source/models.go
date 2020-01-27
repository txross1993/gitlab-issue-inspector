package source

import (
	"time"
)

// User represents a user of GitLab
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar_url"` // Link to user's avatar image
	Link     string `json:"web_url"`    // Link to a user's profile
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
	ID                   int                  `json:"id"`         // Global ID of issue
	IID                  int                  `json:"iid"`        // ID of issue within the project
	Author               User                 `json:"author"`     // Issue creator
	Users                []User               `json:"assignees"`  // Assigned users
	State                string               `json:"state"`      // Open, Closed
	ProjectID            int                  `json:"project_id"` // Glboal project ID
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
	ClosedAt             *time.Time           `json:"closed_at"`
	DueDate              string               `json:"due_date"` // Form of YYYY-MM-DD
	TaskCompletionStatus TaskCompletionStatus `json:"task_completion_status"`
	TimeStats            TimeStats            `json:"time_stats"`
	Link                 string               `json:"web_url"`
	Labels               []string             `json:"labels"`
}

// Project represents the struct for a GitLab project
type Project struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Link           string    `json:"web_url"`
	CreatedAt      time.Time `json:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	Users          []int     // Slice of UserIDs; No json binding for users as this is a separate API request
}

// Notes are the activities that occurred on an issue
type Note struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}
