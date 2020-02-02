package source

import (
	"net/http"
	"sync"
	"time"

	data "github.com/txross1993/gitlab-issue-inspector/data"
)

// Uses fan-out concurrency pattern to fetch all projects and their associated users

// FetchProjects concurrently retrieves all project data provided a list of project IDs
// and returns read-only channels of Project and errors
func FetchProjects(client *http.Client, projectIDs []int) (<-chan data.Project, <-chan error) {
	var wg sync.WaitGroup
	wg.Add(len(projectIDs))
	out := make(chan data.Project)
	errs := make(chan error)

	for _, projectID := range projectIDs {
		go func(projectID int) {
			project, err := fetchProject(client, projectID)
			if err != nil {
				errs <- err
			}
			out <- project
			wg.Done()
		}(projectID)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	return out, errs
}

// FetchUsers concurrently retrieves all user data provided a list of user IDs
// and returns read-only channels of User and errors
func FetchUsers(client *http.Client, userIDs []int) (<-chan data.User, <-chan error) {
	var wg sync.WaitGroup
	wg.Add(len(userIDs))
	out := make(chan data.User)
	errs := make(chan error)

	for _, userID := range userIDs {
		go func(userID int) {
			user, err := fetchUser(client, userID)
			if err != nil {
				errs <- err
			}
			out <- user
			wg.Done()
		}(userID)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	return out, errs
}

type IssueNote struct {
	ProjectID int
	IssueIID  int
}

func FetchNotes(client *http.Client, issueNotes []IssueNote, updatedAt string) (<-chan []data.Note, <-chan error) {
	var wg sync.WaitGroup
	out := make(chan []data.Note)
	errs := make(chan error)

	wg.Add(len(issueNotes))

	for _, issue := range issueNotes {
		go func(issue IssueNote) {
			notes, err := fetchNotes(issue.ProjectID, issue.IssueIID, client)
			if err != nil {
				errs <- err
			}
			notes, err = filterNotesByUpdatedAt(notes, updatedAt)
			if err != nil {
				errs <- err
			}
			out <- notes
		}(issue)
	}

	go func() {
		wg.Wait()
		close(out)
		close(errs)
	}()

	return out, errs
}

func filterNotesByUpdatedAt(notes []data.Note, updatedAt string) ([]data.Note, error) {
	updatedAtTime, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}

	var filteredNotes []data.Note
	for _, note := range notes {
		if note.UpdatedAt.Before(updatedAtTime) {
			break
		}
		filteredNotes = append(filteredNotes, note)
	}

	return filteredNotes, nil
}
