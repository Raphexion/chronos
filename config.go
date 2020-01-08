package main

import (
	"io/ioutil"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var (
	// DefaultURL is the default URL
	DefaultURL = "https://myJira.atlassian.net"
	// DefaultMail is the default mail
	DefaultMail = "myLogin@example.com"
	// DefaultAPIKey is the default API key
	DefaultAPIKey = "1234ABCD"
	// DefaultUsername is the default username
	DefaultUsername = "myUserName"
	// DefaultWeeksLookback number of weeks to look back
	DefaultWeeksLookback = 3
)

// Jira represent all configuration for Jira
type Jira struct {
	URL           string `yaml:"url"`
	Mail          string `yaml:"mail"`
	APIKey        string `yaml:"apikey"`
	Username      string `yaml:"username"`
	WeeksLookback int    `yaml:"weekslookback"`
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

// DefaultConfig returns the default config
func DefaultConfig() (config ChronosConfig) {
	config = ChronosConfig{
		Jira: Jira{
			URL:           DefaultURL,
			Mail:          DefaultMail,
			APIKey:        DefaultAPIKey,
			Username:      DefaultUsername,
			WeeksLookback: DefaultWeeksLookback,
		},
	}
	return
}

// GenerateExampleConfig will write an example configuration to file
func GenerateExampleConfig() error {
	config := DefaultConfig()
	data, err := yaml.Marshal(&config)

	usr, err := user.Current()
	configFile := filepath.Join(usr.HomeDir, "chronos.yaml")
	err = ioutil.WriteFile(configFile, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
