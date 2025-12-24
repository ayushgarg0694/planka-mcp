package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ayushgarg/mcp-planka/internal/planka"
)

// getTools returns the list of available tools
func (s *Server) getTools() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "get_projects",
			"description": "Get all projects",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_project",
			"description": "Get a project by ID",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"projectId": map[string]interface{}{
						"type":        "string",
						"description": "The project ID",
					},
				},
				"required": []string{"projectId"},
			},
		},
		{
			"name":        "create_project",
			"description": "Create a new project",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The project name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "The project description",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			"name":        "get_boards",
			"description": "Get all boards for a project",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"projectId": map[string]interface{}{
						"type":        "string",
						"description": "The project ID",
					},
				},
				"required": []string{"projectId"},
			},
		},
		{
			"name":        "get_board",
			"description": "Get a board by ID",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"boardId": map[string]interface{}{
						"type":        "string",
						"description": "The board ID",
					},
				},
				"required": []string{"boardId"},
			},
		},
		{
			"name":        "create_board",
			"description": "Create a new board",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The board name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "The board description",
					},
					"projectId": map[string]interface{}{
						"type":        "string",
						"description": "The project ID",
					},
				},
				"required": []string{"name", "projectId"},
			},
		},
		{
			"name":        "get_lists",
			"description": "Get all lists for a board",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"boardId": map[string]interface{}{
						"type":        "string",
						"description": "The board ID",
					},
				},
				"required": []string{"boardId"},
			},
		},
		{
			"name":        "get_list",
			"description": "Get a list by ID",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"listId": map[string]interface{}{
						"type":        "string",
						"description": "The list ID",
					},
				},
				"required": []string{"listId"},
			},
		},
		{
			"name":        "create_list",
			"description": "Create a new list",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The list name",
					},
					"boardId": map[string]interface{}{
						"type":        "string",
						"description": "The board ID",
					},
					"position": map[string]interface{}{
						"type":        "number",
						"description": "The list position",
					},
				},
				"required": []string{"name", "boardId"},
			},
		},
		{
			"name":        "get_cards",
			"description": "Get all cards for a list",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"listId": map[string]interface{}{
						"type":        "string",
						"description": "The list ID",
					},
				},
				"required": []string{"listId"},
			},
		},
		{
			"name":        "get_card",
			"description": "Get a card by ID",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "create_card",
			"description": "Create a new card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The card name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "The card description",
					},
					"listId": map[string]interface{}{
						"type":        "string",
						"description": "The list ID",
					},
					"position": map[string]interface{}{
						"type":        "number",
						"description": "The card position",
					},
					"dueDate": map[string]interface{}{
						"type":        "string",
						"description": "The due date (ISO 8601 format)",
					},
				},
				"required": []string{"name", "listId"},
			},
		},
		{
			"name":        "update_card",
			"description": "Update a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The card name",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "The card description",
					},
					"listId": map[string]interface{}{
						"type":        "string",
						"description": "The list ID (to move card)",
					},
					"position": map[string]interface{}{
						"type":        "number",
						"description": "The card position",
					},
					"dueDate": map[string]interface{}{
						"type":        "string",
						"description": "The due date (ISO 8601 format)",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "delete_card",
			"description": "Delete a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "move_card",
			"description": "Move a card to a different list",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
					"listId": map[string]interface{}{
						"type":        "string",
						"description": "The target list ID",
					},
					"position": map[string]interface{}{
						"type":        "number",
						"description": "The card position in the new list",
					},
				},
				"required": []string{"cardId", "listId"},
			},
		},
		{
			"name":        "get_tasks",
			"description": "Get all tasks for a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "create_task",
			"description": "Create a new task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The task name",
					},
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
					"position": map[string]interface{}{
						"type":        "number",
						"description": "The task position",
					},
				},
				"required": []string{"name", "cardId"},
			},
		},
		{
			"name":        "update_task",
			"description": "Update a task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"taskId": map[string]interface{}{
						"type":        "string",
						"description": "The task ID",
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": "The task name",
					},
					"isCompleted": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether the task is completed",
					},
					"position": map[string]interface{}{
						"type":        "number",
						"description": "The task position",
					},
				},
				"required": []string{"taskId"},
			},
		},
		{
			"name":        "delete_task",
			"description": "Delete a task",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"taskId": map[string]interface{}{
						"type":        "string",
						"description": "The task ID",
					},
				},
				"required": []string{"taskId"},
			},
		},
		{
			"name":        "get_comments",
			"description": "Get all comments for a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "create_comment",
			"description": "Create a new comment",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "The comment text",
					},
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"text", "cardId"},
			},
		},
		{
			"name":        "delete_comment",
			"description": "Delete a comment",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"commentId": map[string]interface{}{
						"type":        "string",
						"description": "The comment ID",
					},
				},
				"required": []string{"commentId"},
			},
		},
		{
			"name":        "get_stopwatch",
			"description": "Get the stopwatch for a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "start_stopwatch",
			"description": "Start the stopwatch for a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "stop_stopwatch",
			"description": "Stop the stopwatch for a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
		{
			"name":        "reset_stopwatch",
			"description": "Reset the stopwatch for a card",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"cardId": map[string]interface{}{
						"type":        "string",
						"description": "The card ID",
					},
				},
				"required": []string{"cardId"},
			},
		},
	}
}

// callTool calls a tool by name with the given arguments
func (s *Server) callTool(name string, arguments map[string]interface{}) (string, error) {
	switch name {
	case "get_projects":
		return s.handleGetProjects()
	case "get_project":
		return s.handleGetProject(arguments)
	case "create_project":
		return s.handleCreateProject(arguments)
	case "get_boards":
		return s.handleGetBoards(arguments)
	case "get_board":
		return s.handleGetBoard(arguments)
	case "create_board":
		return s.handleCreateBoard(arguments)
	case "get_lists":
		return s.handleGetLists(arguments)
	case "get_list":
		return s.handleGetList(arguments)
	case "create_list":
		return s.handleCreateList(arguments)
	case "get_cards":
		return s.handleGetCards(arguments)
	case "get_card":
		return s.handleGetCard(arguments)
	case "create_card":
		return s.handleCreateCard(arguments)
	case "update_card":
		return s.handleUpdateCard(arguments)
	case "delete_card":
		return s.handleDeleteCard(arguments)
	case "move_card":
		return s.handleMoveCard(arguments)
	case "get_tasks":
		return s.handleGetTasks(arguments)
	case "create_task":
		return s.handleCreateTask(arguments)
	case "update_task":
		return s.handleUpdateTask(arguments)
	case "delete_task":
		return s.handleDeleteTask(arguments)
	case "get_comments":
		return s.handleGetComments(arguments)
	case "create_comment":
		return s.handleCreateComment(arguments)
	case "delete_comment":
		return s.handleDeleteComment(arguments)
	case "get_stopwatch":
		return s.handleGetStopwatch(arguments)
	case "start_stopwatch":
		return s.handleStartStopwatch(arguments)
	case "stop_stopwatch":
		return s.handleStopStopwatch(arguments)
	case "reset_stopwatch":
		return s.handleResetStopwatch(arguments)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// Helper functions to handle each tool

func (s *Server) handleGetProjects() (string, error) {
	projects, err := s.client.GetProjects()
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetProject(args map[string]interface{}) (string, error) {
	projectID, ok := args["projectId"].(string)
	if !ok {
		return "", fmt.Errorf("missing projectId")
	}
	project, err := s.client.GetProject(projectID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleCreateProject(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("missing name")
	}
	req := planka.CreateProjectRequest{
		Name: name,
	}
	if desc, ok := args["description"].(string); ok {
		req.Description = desc
	}
	project, err := s.client.CreateProject(req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetBoards(args map[string]interface{}) (string, error) {
	projectID, ok := args["projectId"].(string)
	if !ok {
		return "", fmt.Errorf("missing projectId")
	}
	boards, err := s.client.GetBoards(projectID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(boards, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetBoard(args map[string]interface{}) (string, error) {
	boardID, ok := args["boardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing boardId")
	}
	board, err := s.client.GetBoard(boardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(board, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleCreateBoard(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("missing name")
	}
	projectID, ok := args["projectId"].(string)
	if !ok {
		return "", fmt.Errorf("missing projectId")
	}
	req := planka.CreateBoardRequest{
		Name:      name,
		ProjectID: projectID,
	}
	if desc, ok := args["description"].(string); ok {
		req.Description = desc
	}
	board, err := s.client.CreateBoard(req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(board, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetLists(args map[string]interface{}) (string, error) {
	boardID, ok := args["boardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing boardId")
	}
	lists, err := s.client.GetLists(boardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(lists, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetList(args map[string]interface{}) (string, error) {
	listID, ok := args["listId"].(string)
	if !ok {
		return "", fmt.Errorf("missing listId")
	}
	list, err := s.client.GetList(listID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleCreateList(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("missing name")
	}
	boardID, ok := args["boardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing boardId")
	}
	req := planka.CreateListRequest{
		Name:    name,
		BoardID: boardID,
	}
	// Position is required - use provided value or default
	if pos, ok := args["position"].(float64); ok && pos > 0 {
		req.Position = pos
	} else {
		req.Position = 65535 // Default position
	}
	list, err := s.client.CreateList(req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetCards(args map[string]interface{}) (string, error) {
	listID, ok := args["listId"].(string)
	if !ok {
		return "", fmt.Errorf("missing listId")
	}
	cards, err := s.client.GetCards(listID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(cards, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetCard(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	card, err := s.client.GetCard(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleCreateCard(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("missing name")
	}
	listID, ok := args["listId"].(string)
	if !ok {
		return "", fmt.Errorf("missing listId")
	}
	req := planka.CreateCardRequest{
		Name:   name,
		ListID: listID,
	}
	if desc, ok := args["description"].(string); ok {
		req.Description = desc
	}
	if pos, ok := args["position"].(float64); ok {
		req.Position = pos
	}
	if dueDateStr, ok := args["dueDate"].(string); ok {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return "", fmt.Errorf("invalid dueDate format: %w", err)
		}
		req.DueDate = &dueDate
	}
	card, err := s.client.CreateCard(req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleUpdateCard(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	req := planka.UpdateCardRequest{}
	if name, ok := args["name"].(string); ok {
		req.Name = &name
	}
	if desc, ok := args["description"].(string); ok {
		req.Description = &desc
	}
	if listID, ok := args["listId"].(string); ok {
		req.ListID = &listID
	}
	if pos, ok := args["position"].(float64); ok {
		req.Position = &pos
	}
	if dueDateStr, ok := args["dueDate"].(string); ok {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return "", fmt.Errorf("invalid dueDate format: %w", err)
		}
		req.DueDate = &dueDate
	}
	card, err := s.client.UpdateCard(cardID, req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleDeleteCard(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	if err := s.client.DeleteCard(cardID); err != nil {
		return "", err
	}
	return `{"success": true}`, nil
}

func (s *Server) handleMoveCard(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	listID, ok := args["listId"].(string)
	if !ok {
		return "", fmt.Errorf("missing listId")
	}
	position := 0.0
	if pos, ok := args["position"].(float64); ok {
		position = pos
	}
	card, err := s.client.MoveCard(cardID, listID, position)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleGetTasks(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	tasks, err := s.client.GetTasks(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleCreateTask(args map[string]interface{}) (string, error) {
	name, ok := args["name"].(string)
	if !ok {
		return "", fmt.Errorf("missing name")
	}
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	req := planka.CreateTaskRequest{
		Name:   name,
		CardID: cardID,
	}
	if pos, ok := args["position"].(float64); ok {
		req.Position = pos
	}
	task, err := s.client.CreateTask(req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleUpdateTask(args map[string]interface{}) (string, error) {
	taskID, ok := args["taskId"].(string)
	if !ok {
		return "", fmt.Errorf("missing taskId")
	}
	req := planka.UpdateTaskRequest{}
	if name, ok := args["name"].(string); ok {
		req.Name = &name
	}
	if isCompleted, ok := args["isCompleted"].(bool); ok {
		req.IsCompleted = &isCompleted
	}
	if pos, ok := args["position"].(float64); ok {
		req.Position = &pos
	}
	task, err := s.client.UpdateTask(taskID, req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleDeleteTask(args map[string]interface{}) (string, error) {
	taskID, ok := args["taskId"].(string)
	if !ok {
		return "", fmt.Errorf("missing taskId")
	}
	if err := s.client.DeleteTask(taskID); err != nil {
		return "", err
	}
	return `{"success": true}`, nil
}

func (s *Server) handleGetComments(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	comments, err := s.client.GetComments(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleCreateComment(args map[string]interface{}) (string, error) {
	text, ok := args["text"].(string)
	if !ok {
		return "", fmt.Errorf("missing text")
	}
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	req := planka.CreateCommentRequest{
		Text:   text,
		CardID: cardID,
	}
	comment, err := s.client.CreateComment(req)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(comment, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleDeleteComment(args map[string]interface{}) (string, error) {
	commentID, ok := args["commentId"].(string)
	if !ok {
		return "", fmt.Errorf("missing commentId")
	}
	if err := s.client.DeleteComment(commentID); err != nil {
		return "", err
	}
	return `{"success": true}`, nil
}

func (s *Server) handleGetStopwatch(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	stopwatch, err := s.client.GetStopwatch(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(stopwatch, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleStartStopwatch(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	stopwatch, err := s.client.StartStopwatch(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(stopwatch, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleStopStopwatch(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	stopwatch, err := s.client.StopStopwatch(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(stopwatch, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) handleResetStopwatch(args map[string]interface{}) (string, error) {
	cardID, ok := args["cardId"].(string)
	if !ok {
		return "", fmt.Errorf("missing cardId")
	}
	stopwatch, err := s.client.ResetStopwatch(cardID)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(stopwatch, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

