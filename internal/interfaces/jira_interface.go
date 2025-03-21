// Package interfaces provides interfaces for Jira operations
package interfaces

import (
	"context"
	"time"
)

// JiraInterface defines the contract for Jira operations
type JiraInterface interface {
	// GetIssue retrieves detailed information about a specific Jira issue
	GetIssue(ctx context.Context, issueKey string) (*IssueDetails, error)

	// SearchIssues searches for Jira issues using JQL
	SearchIssues(ctx context.Context, jql string) ([]IssueSummary, error)

	// ListSprints retrieves all active and future sprints for a given board
	ListSprints(ctx context.Context, boardID string) ([]SprintInfo, error)

	// CreateIssue creates a new Jira issue
	CreateIssue(ctx context.Context, req CreateIssueRequest) (*CreatedIssue, error)

	// UpdateIssue updates an existing Jira issue
	UpdateIssue(ctx context.Context, req UpdateIssueRequest) error

	// ListStatuses retrieves all available issue statuses for a project
	ListStatuses(ctx context.Context, projectKey string) ([]StatusInfo, error)

	// TransitionIssue transitions an issue through its workflow
	TransitionIssue(ctx context.Context, req TransitionRequest) error
}

// IssueDetails represents detailed information about a Jira issue
type IssueDetails struct {
	Key         string
	Summary     string
	Description string
	Status      string
	Assignee    string
	Priority    string
	Subtasks    []SubtaskInfo
	Transitions []TransitionInfo
}

// IssueSummary represents basic information about a Jira issue
type IssueSummary struct {
	Key         string
	Summary     string
	Status      string
	Assignee    string
	Priority    string
}

// SprintInfo represents information about a Jira sprint
type SprintInfo struct {
	ID        string
	Name      string
	State     string
	StartDate time.Time
	EndDate   time.Time
}

// CreateIssueRequest represents the parameters needed to create a new issue
type CreateIssueRequest struct {
	ProjectKey   string
	Summary      string
	Description  string
	IssueType    string
}

// CreatedIssue represents the response after creating a new issue
type CreatedIssue struct {
	Key  string
	ID   string
	URL  string
}

// UpdateIssueRequest represents the parameters needed to update an issue
type UpdateIssueRequest struct {
	IssueKey     string
	Summary      *string
	Description  *string
}

// StatusInfo represents information about a Jira status
type StatusInfo struct {
	ID          string
	Name        string
	Description string
}

// TransitionRequest represents the parameters needed to transition an issue
type TransitionRequest struct {
	IssueKey     string
	TransitionID string
	Comment      *string
}

// SubtaskInfo represents information about a subtask
type SubtaskInfo struct {
	Key      string
	Summary  string
	Status   string
}

// TransitionInfo represents information about an available transition
type TransitionInfo struct {
	ID   string
	Name string
} 