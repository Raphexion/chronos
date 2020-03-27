package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
)

// A TimeEntry represent a worklog that was entered in JIRA
type TimeEntry struct {
	Issue        string
	Summary      string
	Employee     string
	EmailAddress string
	Date         string
	Hours        float32
	Comment      string
	Week         int
}

type timeEntryPredicate func(TimeEntry) bool

func issueAndWorklogToTimeEntry(issue jira.Issue, worklog jira.WorklogRecord) (entry TimeEntry) {
	entry.Issue = issue.Key
	entry.Summary = issue.Fields.Summary
	entry.Employee = worklog.Author.Name
	entry.EmailAddress = worklog.Author.EmailAddress
	entry.Date = time.Time(*worklog.Created).Format("2006-01-02")
	entry.Hours = float32(worklog.TimeSpentSeconds) / 3600
	entry.Comment = worklog.Comment

	_, entry.Week = time.Time(*worklog.Created).ISOWeek()
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

func extractAllWorklogsForIssues(client *jira.Client, issues []jira.Issue) (timeEntries []TimeEntry) {
	for _, issue := range issues {
		key := issue.Key
		worklog, _, err := client.Issue.GetWorklogs(key)
		if err != nil {
			log.Fatalf("[collector] Unable to extract worklog for issue %s", err)
		}
		for _, worklogRecord := range worklog.Worklogs {
			timeEntries = append(timeEntries, issueAndWorklogToTimeEntry(issue, worklogRecord))
		}
	}
	return
}

// ExtractTimeEntriesFromJira extracts the latest worklogs for a user
func ExtractTimeEntriesFromJira(client *jira.Client, config ChronosConfig) ([]TimeEntry, error) {
	searchOpts := &jira.SearchOptions{
		StartAt:    0,
		MaxResults: 1000,
		Expand:     "worklog",
		Fields:     []string{"key", "summary", "worklog"},
	}

	pastDate := CalcPassedDate(config).Format("2006-01-02")
	log.Printf("[collector] Query from %s for user %s", pastDate, config.Jira.Username)
	searchString := fmt.Sprintf("worklogDate >= %s && worklogAuthor = %s", pastDate, config.Jira.Username)
	issues, _, err := client.Issue.Search(searchString, searchOpts)
	if err != nil {
		log.Fatalf("[collector] Search failed %s", err)
		return []TimeEntry{}, err
	}

	log.Printf("JIRA returned %d items", len(issues))

	// To get all worklogs (more than 20), we need to iterate
	// over each issue and do a new request
	timeEntries := extractAllWorklogsForIssues(client, issues)

	employeeTimeEntries := filterTimeEntries(timeEntries, func(worklog TimeEntry) bool {
		return worklog.Employee == config.Jira.Username || strings.HasPrefix(worklog.EmailAddress, config.Jira.Username)
	})

	recentTimeEntries := filterTimeEntries(employeeTimeEntries, func(worklog TimeEntry) bool {
		return worklog.Date > pastDate
	})

	return recentTimeEntries, nil
}

// CalcPassedDate calculates the date in the passed
func CalcPassedDate(config ChronosConfig) time.Time {
	return CalcPassedDateFrom(time.Now(), config)
}

// CalcPassedDateFrom calucates the date in the passed from given date
func CalcPassedDateFrom(from time.Time, config ChronosConfig) time.Time {
	weekday := int(from.Weekday()) - 1
	days := config.WeeksLookback*7 + weekday

	ret := from.AddDate(0, 0, -days)

	return ret
}
