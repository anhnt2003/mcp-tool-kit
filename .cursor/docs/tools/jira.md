Jira Tool Registration and Handlers Documentation

# Overview

This Go package provides a set of tools to interact with Jira using the mcp-go framework. It registers various Jira-related functionalities to the MCPServer, enabling operations like retrieving, searching, updating, creating, and transitioning Jira issues.

The package utilizes several external dependencies, including go-atlassian for Jira API interaction and dev-kit for additional utilities and services.

## Tools Registration

The RegisterJiraTool function is responsible for registering multiple Jira-related tools with the MCPServer. These tools allow for various Jira operations such as:

1. jira_get_issue

Description: Retrieves detailed information about a specific Jira issue, including status, assignee, description, subtasks, and available transitions.

Parameters:

issue_key (string, required): The unique identifier of the Jira issue (e.g., KP-2, PROJ-123).

2. jira_search_issue

Description: Searches for Jira issues using JQL (Jira Query Language) and returns details such as summary, status, assignee, and priority.

Parameters:

jql (string, required): JQL query string (e.g., project = KP AND status = "In Progress").

3. jira_list_sprints

Description: Lists all active and future sprints for a given Jira board.

Parameters:

board_id (string, required): Numeric ID of the Jira board (can be found in the board URL).

4. jira_create_issue

Description: Creates a new Jira issue and returns the created issue's key, ID, and URL.

Parameters:

project_key (string, required): The project identifier where the issue will be created.

summary (string, required): Title of the issue.

description (string, required): Detailed explanation of the issue.

issue_type (string, required): Type of issue (e.g., Bug, Task, Story, Epic).

5. jira_update_issue

Description: Updates an existing Jira issue with new details. Only provided fields will be updated.

Parameters:

issue_key (string, required): The identifier of the issue to update.

summary (string, optional): New title for the issue.

description (string, optional): New description for the issue.

6. jira_list_statuses

Description: Retrieves all available issue statuses for a specific Jira project.

Parameters:

project_key (string, required): The project identifier.

7. jira_transition_issue

Description: Transitions an issue through its workflow using a valid transition ID.

Parameters:

issue_key (string, required): The issue to transition.

transition_id (string, required): Transition ID.

comment (string, optional): Optional comment to add during the transition.

Handlers Implementation

Each tool has a corresponding handler function that processes the input arguments and interacts with the Jira API.

1. jiraIssueHandler

Retrieves details of a Jira issue, including subtasks and available transitions.

Returns formatted issue information.

2. jiraSearchHandler

Executes a JQL-based search and formats the results.

Displays issue key, summary, status, priority, and other details.

3. jiraListSprintHandler

Fetches sprint details for a specific board.

Returns a list of sprint IDs, names, states, and dates.

4. jiraCreateIssueHandler

Creates a new Jira issue using the provided parameters.

Returns the created issue's key, ID, and URL.

5. jiraUpdateIssueHandler

Updates a Jira issue with new summary or description.

Sends an update request to Jira and confirms the changes.

6. jiraGetStatusesHandler

Retrieves and formats all available statuses for a given project.

Returns a list of issue types with their associated statuses.

7. jiraTransitionIssueHandler

Moves an issue to a new status using the specified transition ID.

Optionally adds a comment during the transition.

### Error Handling

Each handler function implements error handling by:

Checking for required parameters.

Validating data types.

Handling API response errors.

Using timeouts (4 * time.Second) to prevent long-running requests.

#### Conclusion

This package provides a comprehensive set of Jira tools that streamline issue management, sprint tracking, and workflow transitions. By integrating with MCPServer, these tools facilitate automation and efficient task handling within Jira.

