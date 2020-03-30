package main

import (
	"bytes"
	"fmt"
	"sort"
)

// Command represent a low-level presentation command
type Command interface {
	isCommand()
}

type clearWeek struct{}
type clearDate struct{}
type clearIssue struct{}
type newWeek struct{ week int }
type newDate struct{ date string }
type newIssue struct{ issue, summary string }
type summaryDate struct{}
type summaryWeek struct{}
type noteHours struct {
	hours   float32
	comment string
}
type printNewIssue struct{}
type printSameIssue struct{}
type printIssueSummary struct{ issue, summary string }

func (clearWeek) isCommand()         {}
func (clearDate) isCommand()         {}
func (clearIssue) isCommand()        {}
func (newWeek) isCommand()           {}
func (newDate) isCommand()           {}
func (newIssue) isCommand()          {}
func (summaryWeek) isCommand()       {}
func (summaryDate) isCommand()       {}
func (noteHours) isCommand()         {}
func (printNewIssue) isCommand()     {}
func (printSameIssue) isCommand()    {}
func (printIssueSummary) isCommand() {}

// ExtractIssueSummaries will take many timeEntries and extract their summaries
func ExtractIssueSummaries(timeEntries []TimeEntry) (issues []string, summaries map[string]string) {
	issueSet := make(map[string]bool)
	summaries = make(map[string]string)
	for _, entry := range timeEntries {
		issueSet[entry.Issue] = true
		summaries[entry.Issue] = entry.Summary
	}

	for issue := range issueSet {
		issues = append(issues, issue)
	}

	sort.Strings(issues)
	return issues, summaries
}

// BuildCommands will take time entries and create low-level commands.
// We do it like this to avoid subtle bugs
func BuildCommands(timeEntries []TimeEntry) (commands []Command) {
	sort.Slice(timeEntries, func(i, j int) bool {
		// If date is the same, sort on issue
		if timeEntries[i].Date == timeEntries[j].Date {
			return timeEntries[i].Issue < timeEntries[j].Issue
		}

		// However, sort primarily on date
		return timeEntries[i].Date < timeEntries[j].Date
	})

	var currentWeek = 0
	var currentDate = ""
	var currentIssue = ""

	for _, timeEntry := range timeEntries {
		if timeEntry.Week != currentWeek {
			if currentWeek != 0 {
				commands = append(commands, summaryDate{})
				commands = append(commands, summaryWeek{})

				commands = append(commands, clearIssue{})
				commands = append(commands, clearDate{})
				commands = append(commands, clearWeek{})
			}

			commands = append(commands, newWeek{week: timeEntry.Week})
			commands = append(commands, newDate{date: timeEntry.Date})
			commands = append(commands, newIssue{issue: timeEntry.Issue, summary: timeEntry.Summary})

			commands = append(commands, noteHours{hours: timeEntry.Hours, comment: timeEntry.Comment})
			commands = append(commands, printNewIssue{})

			currentWeek = timeEntry.Week
			currentDate = timeEntry.Date
			currentIssue = timeEntry.Issue

		} else if timeEntry.Date != currentDate {
			commands = append(commands, summaryDate{})

			commands = append(commands, clearIssue{})
			commands = append(commands, clearDate{})

			commands = append(commands, newDate{date: timeEntry.Date})
			commands = append(commands, newIssue{issue: timeEntry.Issue, summary: timeEntry.Summary})

			commands = append(commands, noteHours{hours: timeEntry.Hours, comment: timeEntry.Comment})
			commands = append(commands, printNewIssue{})

			currentDate = timeEntry.Date
			currentIssue = timeEntry.Issue

		} else if timeEntry.Issue != currentIssue {
			commands = append(commands, clearIssue{})

			commands = append(commands, newIssue{issue: timeEntry.Issue, summary: timeEntry.Summary})

			commands = append(commands, noteHours{hours: timeEntry.Hours, comment: timeEntry.Comment})
			commands = append(commands, printNewIssue{})

			currentIssue = timeEntry.Issue
		} else {
			commands = append(commands, noteHours{hours: timeEntry.Hours, comment: timeEntry.Comment})
			commands = append(commands, printSameIssue{})
		}
	}

	// Always finish by clearing and showing missing summaries
	if len(commands) > 0 {
		commands = append(commands, clearIssue{})
	}

	if len(commands) > 0 {
		commands = append(commands, summaryDate{})
		commands = append(commands, clearDate{})
	}

	if len(commands) > 0 {
		commands = append(commands, summaryWeek{})
		commands = append(commands, clearWeek{})
	}

	if len(commands) > 0 {
		issues, summaries := ExtractIssueSummaries(timeEntries)

		for _, issue := range issues {
			summary := summaries[issue]
			commands = append(commands, printIssueSummary{issue: issue, summary: summary})
		}
	}

	return
}

// PrettyPrint converts commands to bytes buffer.
// Instead of printing directly to stdout we make
// the code more testable using a bytes.Buffer that
// we can easily inspect
func PrettyPrint(commands []Command) (out bytes.Buffer) {
	showComments := false
	var weekTotal float32 = 0.0
	var dateTotal float32 = 0.0
	var issueTotal float32 = 0.0
	var issueHours float32 = 0.0

	var week int = 0
	var date string = ""
	var issue string = ""
	var issueText string = ""
	var comment string = ""

	for _, command := range commands {
		switch cmd := command.(type) {

		case clearWeek:
			weekTotal = 0.0
		case clearDate:
			dateTotal = 0.0
			issueText = ""
		case clearIssue:
			issueTotal = 0.0
			issueHours = 0.0

		case newWeek:
			week = cmd.week
			out.WriteString("===========================\n")
			out.WriteString(fmt.Sprintf("Week %2d\n", week))
			out.WriteString("===========================\n")
			out.WriteString("\n")

		case newDate:
			date = cmd.date
			out.WriteString(fmt.Sprintf("%s\n", date))

		case newIssue:
			issue = cmd.issue
			issueText = cmd.summary

		case summaryDate:
			if date != "" {
				out.WriteString("\t------------------\n")
				out.WriteString(fmt.Sprintf("\t\t %6.2f\n", dateTotal))
			}

		case summaryWeek:
			if week > 0 {
				out.WriteString("\n")
				out.WriteString(fmt.Sprintf("\tTotal:   %6.2f\n", weekTotal))
				out.WriteString("\n")
			}

		// note down the time for an issue
		case noteHours:
			weekTotal += cmd.hours
			dateTotal += cmd.hours
			issueTotal += cmd.hours
			issueHours = cmd.hours
			comment = cmd.comment

		case printNewIssue:
			if issue != "" {
				if comment != "" && showComments {
					out.WriteString(fmt.Sprintf("\t%s: %6.2f %s // %s\n", issue, issueHours, issueText, comment))
				} else {
					out.WriteString(fmt.Sprintf("\t%s: %6.2f %s\n", issue, issueHours, issueText))
				}
			}

		case printSameIssue:
			if comment != "" && showComments {
				out.WriteString(fmt.Sprintf("\t    \\--: %6.2f // %s\n", issueHours, comment))
			} else {
				out.WriteString(fmt.Sprintf("\t    \\--: %6.2f\n", issueHours))
			}

		case printIssueSummary:
			// issue := cmd.issue
			// summary := cmd.summary
			// TODO: out.WriteString(fmt.Sprintf("%s: %s\n", issue, summary))
		}
	}
	return
}

// PrettyPrintBrief converts commands to bytes buffer.
// Instead of printing directly to stdout we make
// the code more testable using a bytes.Buffer that
// we can easily inspect
func PrettyPrintBrief(commands []Command) (out bytes.Buffer) {
	var weekTotal float32 = 0.0
	var dateTotal float32 = 0.0
	var issueTotal float32 = 0.0

	var week int = 0

	for _, command := range commands {
		switch cmd := command.(type) {

		case clearWeek:
			weekTotal = 0.0
		case clearDate:
			dateTotal = 0.0
		case clearIssue:
			issueTotal = 0.0

		case newWeek:
			week = cmd.week

		case summaryWeek:
			if week > 0 {
				out.WriteString(fmt.Sprintf("Week [%2d]: %6.2f\n", week, weekTotal))
			}

		// note down the time for an issue
		case noteHours:
			weekTotal += cmd.hours
			dateTotal += cmd.hours
			issueTotal += cmd.hours
		}
	}
	return
}

// Print will pretty print the time entries
func Print(timeEntries []TimeEntry) {
	commands := BuildCommands(timeEntries)
	output := PrettyPrint(commands)
	fmt.Printf(output.String())
}

// PrintBrief prints a brief worklog
func PrintBrief(timeEntries []TimeEntry) {
	commands := BuildCommands(timeEntries)
	output := PrettyPrintBrief(commands)
	fmt.Printf(output.String())
}
