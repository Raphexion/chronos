package main

import (
	"flag"
	"log"

	"github.com/andygrunwald/go-jira"
)

var (
	url            = flag.String("url", "", "the path to the Jira instance, e.g, https://myjira.atlassian.net")
	mail           = flag.String("mail", "", "your mail your are using when log-in")
	username       = flag.String("username", "", "username, e.g, nijo")
	apikey         = flag.String("api-key", "", "JIRA api key")
	generateConfig = flag.Bool("generate-config", false, "Generate and example config in home folder")
)

func main() {
	flag.Parse()

	if *generateConfig {
		GenerateExampleConfigInHome()
		return
	}

	var config ChronosConfig
	var err error
	if *url == "" || *mail == "" || *username == "" || *apikey == "" {
		config, err = ReadConfig()
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		config = CommandlineConfig(*url, *mail, *username, *apikey)
	}

	tp := jira.BasicAuthTransport{
		Username: config.Jira.Mail,
		Password: config.Jira.APIKey,
	}

	client, err := jira.NewClient(tp.Client(), config.Jira.URL)
	if err != nil {
		log.Fatal(err)
		return
	}

	timeEntries, nil := ExtractTimeEntriesFromJira(client, config)
	if err != nil {
		log.Fatal(err)
		return
	}

	Print(timeEntries)
}
