package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func logWorkInJIRA(client *jira.Client, config ChronosConfig, issue string, hours, minutes int) error {
	timeString := fmt.Sprintf("%dh %dm", hours, minutes)
	record := &jira.WorklogRecord{
		TimeSpent: timeString,
	}
	_, _, err := client.Issue.AddWorklogRecord(issue, record)
	if err != nil {
		return err
	}

	return nil
}
