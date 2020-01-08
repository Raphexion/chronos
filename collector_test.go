package main

import (
	"testing"
	"time"
)

func TestCalcPassedDate(t *testing.T) {
	config := DefaultConfig()
	config.WeeksLookback = 1

	layoutISO := "2006-01-02"
	mockNow, _ := time.Parse(layoutISO, "2020-01-08")
	correctDate, _ := time.Parse(layoutISO, "2019-12-30")

	date := CalcPassedDateFrom(mockNow, config)
	if date != correctDate {
		t.Errorf("CalcPassedDateFrom is wrong, got: %s, want: %s.", date, correctDate)
	}
}
