package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateAndReadBack(t *testing.T) {
	configFile := filepath.Join(os.TempDir(), "chronos.yaml")

	GenerateExampleConfig(configFile)
	config, err := ReadConfigFile(configFile)

	if err != nil {
		t.Errorf("Unable to read config file %s", err)
	}

	if config.Jira.URL != DefaultURL {
		t.Errorf("Config URL is wrong, got: %s, want: %s.", config.Jira.URL, DefaultURL)
	}
}
