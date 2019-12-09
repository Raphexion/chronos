package main

import (
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
)

// A TimeEntry represent a worklog that was entered in JIRA
type TimeEntry struct {
	Issue    string
	Summary  string
	Employee string
	Date     string
	Hours    float32
	Comment  string
	Week     int
}

type timeEntryPredicate func(TimeEntry) bool

func issueAndWorklogToTimeEntry(issue jira.Issue, worklog jira.WorklogRecord) (entry TimeEntry) {
	entry.Issue = issue.Key
	entry.Summary = issue.Fields.Summary
	entry.Employee = worklog.Author.Name
	entry.Date = time.Time(*worklog.Created).Format("2006-01-02")
	entry.Hours = float32(worklog.TimeSpentSeconds) / 3600
	entry.Comment = worklog.Comment

	_, entry.Week = time.Time(*worklog.Created).ISOWeek()
	return
}

func extractTimeEntriesFromIssues(issues []jira.Issue) (timeEntries []TimeEntry) {
	for _, issue := range issues {
		for _, worklog := range issue.Fields.Worklog.Worklogs {
			timeEntries = append(timeEntries, issueAndWorklogToTimeEntry(issue, worklog))
		}
	}
	return
}

func filterTimeEntries(timeEntries []TimeEntry, predicate timeEntryPredicate) (ret []TimeEntry) {
	for _, record := range timeEntries {
		if predicate(record) {
			ret = append(ret, record)
		}
	}
	return
}

// ExtractTimeEntriesFromJira extracts the latest worklogs for a user
func ExtractTimeEntriesFromJira(client *jira.Client, config ChronosConfig) ([]TimeEntry, error) {
	searchOpts := &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 1000,
		Fields:     []string{"key", "summary", "worklog"},
	}

	pastDate := time.Now().AddDate(0, -3, 0).Format("2006-01-02")
	searchString := fmt.Sprintf("worklogDate >= %s && worklogAuthor = %s", pastDate, config.Jira.Username)
	issues, _, err := client.Issue.Search(searchString, searchOpts)
	if err != nil {
		return []TimeEntry{}, err
	}

	timeEntries := extractTimeEntriesFromIssues(issues)

	employeeTimeEntries := filterTimeEntries(timeEntries, func(worklog TimeEntry) bool {
		return worklog.Employee == config.Jira.Username
	})

	return employeeTimeEntries, nil
}
