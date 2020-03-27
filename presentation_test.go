package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

var issueA = "AA-1234"
var issueB = "AA-1235"
var summaryA = "Summary of issue A"
var summaryB = "Summary of issue B"

var timeEntry1 = TimeEntry{
	Week:     1,
	Date:     "2018-01-01",
	Issue:    issueA,
	Summary:  summaryA,
	Employee: "maxx",
	Hours:    1.0,
	Comment:  "My Comment 111",
}

var timeEntry2 = TimeEntry{
	Week:     1,
	Date:     "2018-01-01",
	Issue:    issueB,
	Summary:  summaryB,
	Employee: "maxx",
	Hours:    2.0,
	Comment:  "My Comment 222",
}

var timeEntry3 = TimeEntry{
	Week:     2,
	Date:     "2018-01-08",
	Issue:    issueB,
	Summary:  summaryB,
	Employee: "maxx",
	Hours:    3.0,
	Comment:  "My Comment 333",
}

var timeEntry4 = TimeEntry{
	Week:     2,
	Date:     "2018-01-09",
	Issue:    issueB,
	Summary:  summaryB,
	Employee: "maxx",
	Hours:    4.0,
	Comment:  "My Comment 444",
}

func TestCommandBuilderEmpty(t *testing.T) {
	commands := BuildCommands([]TimeEntry{})
	amount := len(commands)
	if amount != 0 {
		t.Errorf("Returned wrong number of commands, got: %d, want: %d.", amount, 0)
	}
}

func TestCommandBuilderWithOneEntry(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1})
	amount := len(commands)

	// Expected
	// 1. New week
	// 2. New date
	// 3. New issue
	// 4. Note hours
	// 5. Print new issue
	// 6. Clear issue
	// 7. Summary date
	// 8. Clear date
	// 9. Summary week
	// 10. Clear week
	// 11. Summary of issue

	expected := 11

	if amount != expected {
		t.Errorf("Returned wrong number of commands, got: %d, want: %d.", amount, expected)
		t.Errorf("Commands %+v", commands)
	}

	var total float32 = 0.0
	for _, command := range commands {
		switch cmd := command.(type) {
		case noteHours:
			total += cmd.hours
		}
	}

	if total == 2.0 {
		t.Errorf("Returned wrong amount of hours, got: %f, want: %f.", total, 2.0)
	}
}

func TestCommandBuilderWithTwoEntries(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1, timeEntry1})
	amount := len(commands)

	// Expected
	// 1. New week
	// 2. New date
	// 3. New issue
	// 4. Note hours
	// 5. Print new issue
	// 6. Note hours
	// 7. Print same issue
	// 8. Clear issue
	// 9. Summary date
	// 10. Clear date
	// 11. Summary week
	// 12. Clear week
	// 13. Summary of issue

	expected := 13

	if amount != expected {
		t.Errorf("Returned wrong number of commands, got: %d, want: %d.", amount, expected)
		t.Errorf("Commands %+v", commands)
	}
}

func helpPrettyPrint(commands []Command, goldenFilename string) (output, expected string) {
	outputBuffer := PrettyPrint(commands)
	golden := filepath.Join("testdata", goldenFilename)
	expectedBytes, _ := ioutil.ReadFile(golden)

	output = outputBuffer.String()
	expected = string(expectedBytes)

	// Editors like to remove trailing newlines
	output = strings.TrimRight(output, "\n")
	expected = strings.TrimRight(expected, "\n")

	return
}

func TestPrettyPrintOneEntry(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1})
	output, expected := helpPrettyPrint(commands, "timeEntry1.txt")

	if output != expected {
		t.Errorf("Wrong output, got:\n%s\nexprected:\n%s\n", output, expected)
	}
}

func TestPrettyPrintTwoEntries(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1, timeEntry1})
	output, expected := helpPrettyPrint(commands, "timeEntry11.txt")

	if output != expected {
		t.Errorf("Wrong output, got:\n%s\nexprected:\n%s\n", output, expected)
	}
}

func TestPrettyPrintTwoDifferentEntries(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1, timeEntry2})
	output, expected := helpPrettyPrint(commands, "timeEntry12.txt")

	if output != expected {
		t.Errorf("Wrong output, got:\n%s\nexprected:\n%s\n", output, expected)
	}
}

func TestPrettyPrintTwoDifferentWeeksEntries(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1, timeEntry2, timeEntry3})
	output, expected := helpPrettyPrint(commands, "timeEntry123.txt")

	if output != expected {
		t.Errorf("Wrong output, got:\n%s\nexprected:\n%s\n", output, expected)
	}
}

func TestPrettyPrint1234(t *testing.T) {
	commands := BuildCommands([]TimeEntry{timeEntry1, timeEntry2, timeEntry3, timeEntry4})
	output, expected := helpPrettyPrint(commands, "timeEntry1234.txt")

	if output != expected {
		t.Errorf("Wrong output, got:\n%s\nexprected:\n%s\n", output, expected)
	}
}

func TestExtractIssueSummariesOneEntry(t *testing.T) {
	timeEntries := []TimeEntry{timeEntry1}
	issues, summaries := ExtractIssueSummaries(timeEntries)

	if len(issues) != 1 {
		t.Errorf("Wrong amout of issues in summary, got:\n%d\nexprected:\n%d\n", len(issues), 1)
	}

	summary := summaries[issues[0]]
	expected := summaryA
	if summary != expected {
		t.Errorf("Wrong summary, got: %s, exprected: %s\n", summary, expected)
	}
}

func TestExtractIssueSummariesDoubleIssue(t *testing.T) {
	timeEntries := []TimeEntry{timeEntry1, timeEntry1}
	issues, summaries := ExtractIssueSummaries(timeEntries)

	if len(issues) != 1 {
		t.Errorf("Wrong amout of issues in summary, got: %d, exprected: %d\n", len(issues), 1)
		t.Errorf("%+v", issues)
	}

	summary := summaries[issues[0]]
	expected := summaryA
	if summary != expected {
		t.Errorf("Wrong summary, got: %s, exprected: %s\n", summary, expected)
	}
}

func TestExtractIssueSummariesTwoIssues(t *testing.T) {
	timeEntries := []TimeEntry{timeEntry1, timeEntry2}
	issues, summaries := ExtractIssueSummaries(timeEntries)

	if len(issues) != 2 {
		t.Errorf("Wrong amout of issues in summary, got: %d, exprected: %d\n", len(issues), 2)
		t.Errorf("%+v", issues)
	}

	summaryA := summaries[issues[0]]
	expectedA := summaryA
	if summaryA != expectedA {
		t.Errorf("Wrong summary, got: %s, exprected: %s\n", summaryA, expectedA)
	}

	summaryB := summaries[issues[1]]
	expectedB := summaryB
	if summaryB != expectedB {
		t.Errorf("Wrong summary, got: %s, exprected: %s\n", summaryB, expectedB)
	}
}
