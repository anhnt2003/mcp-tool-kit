package services

import (
	"fmt"
	"log"
	"os"

	"github.com/andygrunwald/go-jira"
)

func InitializeAtlassianClient() (*jira.Client) {
	email := os.Getenv("JIRA_EMAIL")
	apiKey := os.Getenv("JIRA_API_KEY")
	jiraURL := os.Getenv("JIRA_URL")

	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiKey,
	}

	client, err := jira.NewClient(tp.Client(), jiraURL)
	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	fmt.Println("Jira client created successfully")

	return client
}
