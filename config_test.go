package main

import (
	"testing"
)

func TestGenerateAndReadBack(t *testing.T) {
	GenerateExampleConfig()
	config, err := ReadConfig()

	if err != nil {
		t.Errorf("Unable to read config file %s", err)
	}

	if config.Jira.URL != DefaultURL {
		t.Errorf("Config URL is wrong, got: %s, want: %s.", config.Jira.URL, DefaultURL)
	}
}
