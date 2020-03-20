package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/andygrunwald/go-jira"
)

var (
	url            = flag.String("url", "", "the path to the Jira instance, e.g, https://myjira.atlassian.net")
	mail           = flag.String("mail", "", "your mail your are using when log-in")
	username       = flag.String("username", "", "username, e.g, nijo")
	apikey         = flag.String("api-key", "", "JIRA api key")
	generateConfig = flag.Bool("generate-config", false, "generate and example config in home folder")
	logWork        = flag.Bool("logwork", false, "log time in JIRA")
	issue          = flag.String("issue", "", "issue to query or manipulate")
	hours          = flag.Int("hours", 0, "hours to log time")
	minutes        = flag.Int("minutes", 0, "minutes to log time")
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

	if *logWork {
		if *issue != "" && (*hours > 0 || *minutes > 0) {
			err := logWorkInJIRA(client, config, *issue, *hours, *minutes)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("Successfully logged %dh %dm to %s\n", *hours, *minutes, *issue)
			}
		} else {
			log.Fatalf("Unable to log work, need --issue, --hours and/or --minutes")
		}
		return
	}

	timeEntries, nil := ExtractTimeEntriesFromJira(client, config)
	if err != nil {
		log.Fatal(err)
		return
	}

	Print(timeEntries)
}
