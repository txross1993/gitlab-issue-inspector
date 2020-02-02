package source

import (
	"time"
)

func GetModels() []interface{} {
	var models = []interface{}{
		&User{},
		&TimeStat{},
		&TaskCompletionStat{},
		&Issue{},
		&Project{},
		&Note{},
	}
	return models

}

// User represents a user of GitLab
type User struct {
	ID       int    `json:"id" gorm:"PRIMARY_KEY"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar_url"` // Link to user's avatar image
	Link     string `json:"web_url"`    // Link to a user's profile
}

// TimeStat represents the nested structure of an Issue related to time estimation and time spent toward and issue
type TimeStat struct {
	ID             int `json:"-" gorm:"PRIMARY_KEY"`
	TimeEstimate   int `json:"time_estimate"`
	TotalTimeSpent int `json:"total_time_spent"`
}

// TaskCompletionStat represents the nested structure of an Issue related to number of tasks
type TaskCompletionStat struct {
	ID             int `json:"-" gorm:"PRIMARY_KEY"`
	Count          int `json:"count"`
	CompletedCount int `json:"completed_count"`
}

// Issue represents the struct for a GitLab issue
type Issue struct {
	ID                   int                `json:"id" gorm:"PRIMARY_KEY"` // Global ID of issue
	IID                  int                `json:"iid" sql:"index"`       // ID of issue within the project
	AuthorID             int                `json:"-"`
	Author               User               `json:"author" gorm:"foreignkey:AuthorID"` // Issue creator
	UserIDs              []int              `json:"-" gorm:"type:INT[][]"`             // Must be generated internally, not part of JSON response
	Users                []User             `json:"assignees" gorm:"-"`                // Assigned users
	State                string             `json:"state"`                             // Open, Closed
	ProjectID            int                `json:"project_id"`                        // Glboal project ID
	Project              Project            `json:"-" gorm:"foreignkey:ProjectID"`
	CreatedAt            time.Time          `json:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at"`
	ClosedAt             *time.Time         `json:"closed_at"`
	DueDate              string             `json:"due_date"` // Form of YYYY-MM-DD
	TaskCompletionStatID int                `json:"-"`        // Must be generated internally, not part of JSON response
	TaskCompletionStat   TaskCompletionStat `json:"task_completion_status" gorm:"foreignkey:TaskCompletionStatusID"`
	TimeStatID           int                `json:"-"` // Must be generated internally, not part of JSON response
	TimeStat             TimeStat           `json:"time_stats"  gorm:"foreignkey:TimeStatID"`
	Link                 string             `json:"web_url"`
	Labels               []string           `json:"labels" gorm:"type:TEXT[][]|"`
}

func (i *Issue) GetForeignKeyMapping() map[string]string {
	var relationships = make(map[string]string, 3)

	// Add FK for User
	relationships["author_id"] = "users(id)"

	// Add FK for Time Stats
	relationships["time_stats_id"] = "time_stats(id)"

	// Add FK for Time Stats
	relationships["task_completion_stat_id"] = "task_completion_stats(id)"

	return relationships

}

// Project represents the struct for a GitLab project
type Project struct {
	ID             int       `json:"id" gorm:"PRIMARY_KEY"`
	Name           string    `json:"name"`
	Link           string    `json:"web_url"`
	CreatedAt      time.Time `json:"created_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
	UserIDs        []int     `json:"-"` // Must be generated internally, not part of JSON response
}

// Notes are the activities that occurred on an issue
type Note struct {
	ID        int       `json:"id" gorm:"PRIMARY_KEY"`
	IssueID   int       `json:"-"` // Must be generated internally, not part of JSON response
	Issue     Issue     `json:"-" gorm:"foreignkey:IssueID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}

func (n *Note) GetForeignKeyMapping() map[string]string {
	var relationships = make(map[string]string, 1)

	// Add FK for Issue
	relationships["issue_id"] = "issues(id)"

	return relationships

}
