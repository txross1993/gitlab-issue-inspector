-- Assignees
CREATE TABLE IF NOT EXISTS `gitlab_tracking.Assignees` (
    id INT64 NOT NULL, -- Assignee ID in gitlab
    `name` STRING, -- display name in GitLab
    username STRING, -- username in GitLab
    avatar STRING -- URL to avatar.png in GitLab
)

-- Projects contain active projects with related issues
CREATE TABLE IF NOT EXISTS `gitlab_tracking.Projects` (
    id INT64 NOT NULL, -- Project ID in GitLab
    `name` STRING, -- Human readable name of a project
    link STRING, -- URL of the project
    assignees ARRAY<INT64> -- Assignee IDs 
)

-- Issues contains a catalog of known issues
CREATE TABLE IF NOT EXISTS `gitlab_tracking.Issues` (
    id INT64 NOT NULL, -- Global issue id in GitLab
    iid INT64 NOT NULL, -- Project-scoped Issue ID by in GitLab
    project_id INT64, -- Reference project ID
    created_at timestamp NOT NULL, -- When the issue was created
    closed_at timestamp, -- When the issue was closed
    assignees ARRAY<INT64>, -- Assignee IDs
    `state` STRING, -- Open or Closed
    due_date Date, -- When the issue is due
    task_completion_status FLOAT64, -- % of tasks completed to tasks specified
    time_estimate INT64, -- estimated time for task in seconds 
    time_spent INT64 -- time spent on task in seconds
)

-- Logging new events of work!
CREATE TABLE IF NOT EXISTS `gitlab_tracking.Log` (
    event_type STRING, -- one of `issue_created`, `issue_closed`, `time_spent`, `issue_updated`
    timestamp timestamp, -- when the work was done
    issue_id INT64, -- Global issue id in GitLab
    assignee_id INT64, -- who did the work 
    amount INT64 -- In the event of time spent, how much time spent?
)
PARTITION BY DATE(timestamp)
CLUSTER BY assignee_id
OPTIONS (
    partition_expiration_days=60,
    description="A table of work logged clustered by assignee, partitioned by timestamp"
)