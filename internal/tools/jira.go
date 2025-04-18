package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/anhnt2003/mcp-tool-kit/internal/services"

	"github.com/andygrunwald/go-jira"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Jira defines the interface for Jira operations
type JiraInterface interface {
	// GetIssue retrieves detailed information about a specific Jira issue
	GetIssue(ctx context.Context, issueKey string) (*jira.Issue, error)
	
	// CreateIssue creates a new Jira issue
	CreateIssue(ctx context.Context, project string, issueType string, summary string, description string) (*jira.Issue, error)
	
	// UpdateIssue updates an existing Jira issue
	UpdateIssue(ctx context.Context, issueKey string, fields map[string]interface{}) (*jira.Issue, error)
	
	// AssignIssue assigns a Jira issue to a user
	AssignIssue(ctx context.Context, issueKey string, assignee string) error
	
	// TransitionIssue transitions an issue through its workflow
	TransitionIssue(ctx context.Context, issueKey string, transitionID string) error
	
	// AddComment adds a comment to a Jira issue
	AddComment(ctx context.Context, issueKey string, body string) (*jira.Comment, error)
	
	// SearchIssues searches for Jira issues using JQL
	SearchIssues(ctx context.Context, jql string, options *jira.SearchOptions) ([]jira.Issue, error)
}

// jiraServiceImpl implements the JiraService interface
type jiraServiceImpl struct {
	client *jira.Client
}

// GetIssue implements Jira.
func (j *jiraServiceImpl) GetIssue(ctx context.Context, issueKey string) (*jira.Issue, error) {
	issue, _, err := j.client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

// CreateIssue implements Jira.
func (j *jiraServiceImpl) CreateIssue(ctx context.Context, project string, issueType string, summary string, description string) (*jira.Issue, error) {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: description,
			Type: jira.IssueType{
				Name: issueType,
			},
			Project: jira.Project{
				Key: project,
			},
			Summary: summary,
		},
	}
	
	issue, _, err := j.client.Issue.Create(&i)
	if err != nil {
		return nil, err
	}
	return issue, nil
}

// UpdateIssue implements Jira.
func (j *jiraServiceImpl) UpdateIssue(ctx context.Context, issueKey string, fields map[string]interface{}) (*jira.Issue, error) {
	// Use the UpdateIssue method which takes a jiraID and a map of fields
	_, err := j.client.Issue.UpdateIssue(issueKey, fields)
	if err != nil {
		return nil, err
	}
	
	// Get the full issue after update
	return j.GetIssue(ctx, issueKey)
}

// AssignIssue implements Jira.
func (j *jiraServiceImpl) AssignIssue(ctx context.Context, issueKey string, assignee string) error {
	_, err := j.client.Issue.UpdateAssignee(issueKey, &jira.User{
		Name: assignee,
	})
	return err
}

// TransitionIssue implements Jira.
func (j *jiraServiceImpl) TransitionIssue(ctx context.Context, issueKey string, transitionID string) error {
	_, err := j.client.Issue.DoTransition(issueKey, transitionID)
	return err
}

// AddComment implements Jira.
func (j *jiraServiceImpl) AddComment(ctx context.Context, issueKey string, body string) (*jira.Comment, error) {
	comment := &jira.Comment{
		Body: body,
	}
	
	newComment, _, err := j.client.Issue.AddComment(issueKey, comment)
	if err != nil {
		return nil, err
	}
	
	return newComment, nil
}

// SearchIssues implements Jira.
func (j *jiraServiceImpl) SearchIssues(ctx context.Context, jql string, options *jira.SearchOptions) ([]jira.Issue, error) {
	issues, _, err := j.client.Issue.Search(jql, options)
	if err != nil {
		return nil, err
	}
	
	return issues, nil
}

// NewJiraTool creates a new instance of JiraTool
func NewJiraTool(server *server.MCPServer) JiraInterface {
	// Create a new jira service implementation
	jiraTool := &jiraServiceImpl{
		client: services.InitializeAtlassianClient(),
	}
	
	// Add the Jira tool to the MCP server
	if server != nil {
		// Register Jira commands with the MCP server
		
		// Register tool for getting a Jira issue
		getIssueTool := mcp.NewTool("jira_get_issue",
			mcp.WithDescription("Retrieve a Jira issue by its key"),
			mcp.WithString("issue_key",
				mcp.Required(),
				mcp.Description("The key of the Jira issue to retrieve (e.g., PROJ-123)"),
			),
		)
		
		server.AddTool(getIssueTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			issueKey, ok := request.Params.Arguments["issue_key"].(string)
			if !ok {
				return mcp.NewToolResultError("issue_key must be a string"), nil
			}
			
			issue, err := jiraTool.GetIssue(ctx, issueKey)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Convert the issue to a formatted text result
			return mcp.NewToolResultText(
				fmt.Sprintf(
					`Key: %s
					Summary: %s
					Status: %s
					Type: %s
					Assignee: %s
					Priority: %s
					Created: %v
					Updated: %v
					Reporter: %s
					Description: %s	
					Subtasks: %v
					`, 
					issue.Key, 
					issue.Fields.Summary, 
					issue.Fields.Status.Name,
					issue.Fields.Type.Name,
					issue.Fields.Assignee.DisplayName,
					issue.Fields.Priority.Name,
					issue.Fields.Created,
					issue.Fields.Updated,
					issue.Fields.Reporter.DisplayName,
					issue.Fields.Description,
					issue.Fields.Subtasks,
				),
			), nil
		})
		
		// Register tool for creating a Jira issue
		createIssueTool := mcp.NewTool("jira_create_issue",
			mcp.WithDescription("Create a new Jira issue"),
			mcp.WithString("project",
				mcp.Required(),
				mcp.Description("The project key where the issue will be created"),
			),
			mcp.WithString("issue_type",
				mcp.Required(),
				mcp.Description("The type of issue to create (e.g., Bug, Task, Story)"),
			),
			mcp.WithString("summary",
				mcp.Required(),
				mcp.Description("The summary or title of the issue"),
			),
			mcp.WithString("description",
				mcp.Required(),
				mcp.Description("The detailed description of the issue"),
			),
		)
		
		server.AddTool(createIssueTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			project, ok := request.Params.Arguments["project"].(string)
			if !ok {
				return mcp.NewToolResultError("project must be a string"), nil
			}
			
			issueType, ok := request.Params.Arguments["issue_type"].(string)
			if !ok {
				return mcp.NewToolResultError("issue_type must be a string"), nil
			}
			
			summary, ok := request.Params.Arguments["summary"].(string)
			if !ok {
				return mcp.NewToolResultError("summary must be a string"), nil
			}
			
			description, ok := request.Params.Arguments["description"].(string)
			if !ok {
				return mcp.NewToolResultError("description must be a string"), nil
			}
			
			issue, err := jiraTool.CreateIssue(ctx, project, issueType, summary, description)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Return the created issue key and summary as text
			return mcp.NewToolResultText(fmt.Sprintf("Created issue %s: %s", issue.Key, issue.Fields.Summary)), nil
		})
		
		// Register tool for searching Jira issues with JQL
		searchIssuesTool := mcp.NewTool("jira_search_issues",
			mcp.WithDescription("Search for Jira issues using JQL"),
			mcp.WithString("jql",
				mcp.Required(),
				mcp.Description("The JQL query to search for issues"),
			),
			mcp.WithNumber("max_results",
				mcp.Description("Maximum number of results to return"),
			),
		)
		
		server.AddTool(searchIssuesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			jql, ok := request.Params.Arguments["jql"].(string)
			if !ok {
				return mcp.NewToolResultError("jql must be a string"), nil
			}
			
			var maxResults int
			if maxResultsArg, exists := request.Params.Arguments["max_results"]; exists {
				if maxResultsFloat, ok := maxResultsArg.(float64); ok {
					maxResults = int(maxResultsFloat)
				}
			}
			
			searchOptions := &jira.SearchOptions{
				MaxResults: maxResults,
			}
			
			issues, err := jiraTool.SearchIssues(ctx, jql, searchOptions)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Format search results as text
			var resultText strings.Builder
			resultText.WriteString(fmt.Sprintf("Found %d issues:\n\n", len(issues)))
			
			for _, issue := range issues {
				resultText.WriteString(
					fmt.Sprintf(
						`Key: %s
						Summary: %s
						Status: %s
						Type: %s
						Assignee: %s
						Priority: %s
						Created: %v
						Updated: %v
						Reporter: %s
						Description: %s
						Subtasks: %v
						`, 
						issue.Key, 
						issue.Fields.Summary, 
					    issue.Fields.Status.Name,
						issue.Fields.Type.Name,
						issue.Fields.Assignee.DisplayName,
						issue.Fields.Priority.Name,
						issue.Fields.Created,
						issue.Fields.Updated,
						issue.Fields.Reporter.DisplayName,
						issue.Fields.Description,
						issue.Fields.Subtasks,
					),
				)
			}
			
			return mcp.NewToolResultText(resultText.String()), nil
		})
		
		// Register tool for adding comments to issues
		addCommentTool := mcp.NewTool("jira_add_comment",
			mcp.WithDescription("Add a comment to a Jira issue"),
			mcp.WithString("issue_key",
				mcp.Required(),
				mcp.Description("The key of the issue to add a comment to"),
			),
			mcp.WithString("comment",
				mcp.Required(),
				mcp.Description("The comment text to add"),
			),
		)
		
		server.AddTool(addCommentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			issueKey, ok := request.Params.Arguments["issue_key"].(string)
			if !ok {
				return mcp.NewToolResultError("issue_key must be a string"), nil
			}
			
			comment, ok := request.Params.Arguments["comment"].(string)
			if !ok {
				return mcp.NewToolResultError("comment must be a string"), nil
			}
			
			newComment, err := jiraTool.AddComment(ctx, issueKey, comment)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Return a success message with the comment details
			return mcp.NewToolResultText(fmt.Sprintf("Added comment to issue %s: %s", issueKey, newComment.Body)), nil
		})
		
		// Add more commands as needed
		
		// Register tool for transitioning issues (changing status)
		transitionIssueTool := mcp.NewTool("jira_transition_issue",
			mcp.WithDescription("Change the status of a Jira issue"),
			mcp.WithString("issue_key",
				mcp.Required(),
				mcp.Description("The key of the issue to transition"),
			),
			mcp.WithString("transition_id",
				mcp.Required(),
				mcp.Description("The ID of the transition to perform"),
			),
		)
		
		server.AddTool(transitionIssueTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			issueKey, ok := request.Params.Arguments["issue_key"].(string)
			if !ok {
				return mcp.NewToolResultError("issue_key must be a string"), nil
			}
			
			transitionID, ok := request.Params.Arguments["transition_id"].(string)
			if !ok {
				return mcp.NewToolResultError("transition_id must be a string"), nil
			}
			
			err := jiraTool.TransitionIssue(ctx, issueKey, transitionID)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Return a success message
			return mcp.NewToolResultText(fmt.Sprintf("Successfully transitioned issue %s", issueKey)), nil
		})
	}
	
	return jiraTool
}

