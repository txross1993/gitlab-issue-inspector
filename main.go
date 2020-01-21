package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/txross1993/gitlab-issue-inspector/issues"
	"github.com/txross1993/gitlab-issue-inspector/sink"
)

func main() {
	labels := flag.String("labels", "", "Comma-separated list of GitLab issue labels")
	// flag options
	flag.Parse()

	godotenv.Load()

	client := &http.Client{}

	// Fetch the max updatedAt time stored in DB
	sinkClient, err := sink.NewSinkClient()
	if err != nil {
		log.Fatal(err)
	}

	updatedAt, err := sinkClient.GetLastUpdatedTime()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(updatedAt)

	// Fetch issues provided labels
	got, err := issues.Fetch(client, *labels, updatedAt)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(got)

}

// {
// 	// Access Token: [redacted]
// 	"issue_title": "test"
// 	"assignee":
// 	"project":
// 	"project_link":
// 	"due_date":
// 	"created_at":
// 	"closed_at":
// 	"time_estimate":  //seconds
// 	"total_time_spent":  //seconds
// 	"task_completion_status"."count"
// 	"task_completion_status"."completed_count"

// 	/*
// 		If you track the following items:
// 			Issue created_at
// 			Issue due_date
// 			Issue closed_at
// 			Issue human_time_estimate
// 			Issue human_total_time_spent

// 		You can determine the following pieces of information:
// 		- How much more time did you spend on the task than you estimated? This is an indicator for task weight and can be indicative of underestimating, overestimating, or poor story formation.
// 		- How far after/before the due_date was the issue closed? This can indicate the same problems as above, or be a false positive (i.e. closing an issue that has become irrelevant)
// 		- No, we can't use closed_at as an indicator due to side effects. Estimates and time spent are the only true inputs. Though we should be looking for issues to close before or on their due dates.
// 		- Pull in the reference project name
// 	*/

// 	/*
// 	2020-01-20T21:28:05.261Z"
// 	2020-01-20T21:28:05.289Z
// 		API would

// 		Get MAX(updated_at) date
// 		If null, Get all issues
// 		Else, get all issues with label=$LABELS updated_after(MAX(updated_at))

// 		Put delta data updates to DB
// 		{issue_id, issue_title, assignee, created_at, updated_at, closed_at, time_estimate, total_time_spent, task_completion_status}

// 		Do this on schedule - every 1 hr, 10 hrs/day
// 			Cost on this ?

// 		Reporting
// 		- How much time was spent per project per day by assignee
// 		- How does that project relate back to an SR ?
// 		- Can't be the only means of reporting unless everything you do has a project

// 		# ---------------------------------------------------------------------------------
// 		Logging Work
// 		Work is logged when tasks are completed or issues are closed

// 		When Work is logged
// 		Assignee-Issue
// 		{Date, Issue, Task Completed: "task description"}

// 		# ---------------------------------------------------------------------------------
// 		Logging Hours
// 		Hours are logged when time is spent.

// 		When hours are logged
// 		{Date: , Assignee: ,Issue: , Time Spent: }

// 		Issue notes will contain tasks that can be parsed out per user and determine exactly when the task is marked as completed and related back to the parent issue.

// 		notes.body "marked the task * as completed", created_at

// 		The author of the notes change gets credit for the work.
// 	*/

// }
