package services

import (
	"log"
	"os"

	"github.com/andygrunwald/go-jira"
)

func InitializeAtlassianClient() *jira.Client {
	username := os.Getenv("JIRA_USERNAME")
	token := os.Getenv("JIRA_TOKEN")
	jiraURL := os.Getenv("JIRA_URL")

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(tp.Client(), jiraURL)
	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	return client
}
