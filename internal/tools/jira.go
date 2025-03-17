package tools

import (
	"context"
	"fmt"
	"mcp-tool-kit/internal/services"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewJiraTool creates a new instance of JiraTool
func NewJiraTool(server *server.MCPServer) {
	
	// Add a new tool to the server
	jiraGetIssue := mcp.NewTool(
		"jira-get-issue", 
		mcp.WithDescription("Get a Jira issue by ID"), 
		mcp.WithString("issue_key", mcp.Required(), mcp.Description("The unique identifier of the Jira issue (e.g., KP-2, PROJ-123)")),
	)

	server.AddTool(jiraGetIssue, JiraGetIssue)
}

// JiraGetIssue is a tool that gets a Jira issue by ID
func JiraGetIssue(ctx context.Context, req mcp.CallToolRequest) (result *mcp.CallToolResult, err error) {
	client := services.InitializeAtlassianClient()

	issueKey := "KP-2"

	issue, _, err := client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, err
	}
	
	formattedText := fmt.Sprintf("Issue: %s\nSummary: %s\nDescription: %s\nStatus: %s\nAssignee: %s\nCreated: %s\nUpdated: %s",
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.Description,
		issue.Fields.Status.Name,
		issue.Fields.Assignee.Name,
		issue.Fields.Created.String(),
		issue.Fields.Updated.String(),
	)

	result = mcp.NewToolResultText(formattedText)
	return result, nil
}
