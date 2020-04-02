package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/andygrunwald/go-jira"
)

// SprintIssue represent a issue in the sprint
type SprintIssue struct {
	issue    string
	summary  string
	assignee string
}

const unassignedIssue = "Unassigned"

func sprintIssueFromJiraIssue(issue jira.Issue) (ret SprintIssue) {
	ret.issue = issue.Key
	ret.summary = issue.Fields.Summary
	ret.assignee = unassignedIssue
	if issue.Fields.Assignee != nil {
		ret.assignee = issue.Fields.Assignee.EmailAddress
	}
	return
}

func jiraIssuesToSprintIssues(jiraIssues []jira.Issue) (issues []SprintIssue) {
	for _, issue := range jiraIssues {
		issues = append(issues, sprintIssueFromJiraIssue(issue))
	}
	return
}

func unassigned(issue SprintIssue) bool {
	return issue.assignee == unassignedIssue
}

func usersIssue(issue SprintIssue, config ChronosConfig) bool {
	return issue.assignee == config.Jira.Username || strings.HasPrefix(issue.assignee, config.Jira.Username)
}

func keepUsersAndUnassignedIssues(sprintIssues []SprintIssue, config ChronosConfig) (ret []SprintIssue) {
	for _, issue := range sprintIssues {
		if unassigned(issue) || usersIssue(issue, config) {
			ret = append(ret, issue)
		}
	}
	return
}

// UsersIssuesInOpenSprints returns all the user's issues in the open sprints
func UsersIssuesInOpenSprints(client *jira.Client, config ChronosConfig) ([]SprintIssue, error) {
	searchOpts := &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 1000,
		Expand:     "worklog",
		Fields:     []string{"key", "summary", "worklog", "assignee"},
	}

	searchStringTemplate := "resolution = Unresolved AND sprint in openSprints()"
	searchString := fmt.Sprintf(searchStringTemplate)
	jiraIssues, _, err := client.Issue.Search(searchString, searchOpts)
	if err != nil {
		log.Fatalf("[sprint] Search failed %s", err)
		return []SprintIssue{}, err
	}

	allIssues := jiraIssuesToSprintIssues(jiraIssues)
	issues := keepUsersAndUnassignedIssues(allIssues, config)

	return issues, nil
}
