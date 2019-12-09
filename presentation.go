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
type newIssue struct{ issue string }
type summaryDate struct{}
type summaryWeek struct{}
type noteHours struct{ hours float32 }
type printNewIssue struct{}
type printSameIssue struct{}

func (clearWeek) isCommand()      {}
func (clearDate) isCommand()      {}
func (clearIssue) isCommand()     {}
func (newWeek) isCommand()        {}
func (newDate) isCommand()        {}
func (newIssue) isCommand()       {}
func (summaryWeek) isCommand()    {}
func (summaryDate) isCommand()    {}
func (noteHours) isCommand()      {}
func (printNewIssue) isCommand()  {}
func (printSameIssue) isCommand() {}

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
			commands = append(commands, newIssue{issue: timeEntry.Issue})

			commands = append(commands, noteHours{hours: timeEntry.Hours})
			commands = append(commands, printNewIssue{})

			currentWeek = timeEntry.Week
			currentDate = timeEntry.Date
			currentIssue = timeEntry.Issue

		} else if timeEntry.Date != currentDate {
			commands = append(commands, summaryDate{})

			commands = append(commands, clearIssue{})
			commands = append(commands, clearDate{})

			commands = append(commands, newDate{date: timeEntry.Date})
			commands = append(commands, newIssue{issue: timeEntry.Issue})

			commands = append(commands, noteHours{hours: timeEntry.Hours})
			commands = append(commands, printNewIssue{})

			currentDate = timeEntry.Date
			currentIssue = timeEntry.Issue

		} else if timeEntry.Issue != currentIssue {
			commands = append(commands, clearIssue{})

			commands = append(commands, newIssue{issue: timeEntry.Issue})

			commands = append(commands, noteHours{hours: timeEntry.Hours})
			commands = append(commands, printNewIssue{})

			currentIssue = timeEntry.Issue
		} else {
			commands = append(commands, noteHours{hours: timeEntry.Hours})
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

	return
}

// PrettyPrint converts commands to bytes buffer.
// Instead of printing directly to stdout we make
// the code more testable using a bytes.Buffer that
// we can easily inspect
func PrettyPrint(commands []Command) (out bytes.Buffer) {
	var weekTotal float32 = 0.0
	var dateTotal float32 = 0.0
	var issueTotal float32 = 0.0
	var issueHours float32 = 0.0

	var week int = 0
	var date string = ""
	var issue string = ""

	for _, command := range commands {
		switch cmd := command.(type) {

		case clearWeek:
			weekTotal = 0.0
		case clearDate:
			dateTotal = 0.0
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

		case printNewIssue:
			if issue != "" {
				out.WriteString(fmt.Sprintf("\t%s: %6.2f\n", issue, issueHours))
			}

		case printSameIssue:
			out.WriteString(fmt.Sprintf("\t    \\--: %6.2f\n", issueHours))
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
