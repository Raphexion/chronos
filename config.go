package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Jira represent all configuration for Jira
type Jira struct {
	URL      string `yaml:"url"`
	Mail     string `yaml:"mail"`
	APIKey   string `yaml:"apikey"`
	Username string `yaml:"username"`
}

// A ChronosConfig represents all the information we need to
// connect to the JIRA Instance
type ChronosConfig struct {
	Jira
}

// ReadConfig reads a YAML configuration from the home folder
func ReadConfig() (ChronosConfig, error) {
	var config ChronosConfig

	usr, err := user.Current()
	configFile := filepath.Join(usr.HomeDir, "chronos.yaml")
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(raw, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// CommandlineConfig creates a config from user provided information
func CommandlineConfig(url, mail, username, apikey string) (c ChronosConfig) {
	c.Jira.URL = url
	c.Jira.Mail = mail
	c.Jira.Username = username
	c.Jira.APIKey = apikey
	return
}

// GenerateExampleConfig will write an example configuration to file
func GenerateExampleConfig() error {
	config := ChronosConfig{
		Jira: Jira{
			URL:      "https://myJira.atlassian.net",
			Mail:     "myLogin@example.com",
			APIKey:   "1234ABCD",
			Username: "myUserName",
		},
	}

	data, err := yaml.Marshal(&config)

	usr, err := user.Current()
	configFile := filepath.Join(usr.HomeDir, "chronos.yaml")
	err = ioutil.WriteFile(configFile, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
