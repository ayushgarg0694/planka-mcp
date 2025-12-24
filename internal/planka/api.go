package planka

import (
	"encoding/json"
	"fmt"
	"time"
)

// APIResponse represents a generic API response with items
type APIResponse struct {
	Items    interface{} `json:"items"`
	Included interface{} `json:"included,omitempty"`
}

// extractItems extracts items from an API response and unmarshals them into the target type
func extractItems[T any](resp APIResponse) ([]T, error) {
	itemsJSON, err := json.Marshal(resp.Items)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal items: %w", err)
	}
	
	var items []T
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}
	
	return items, nil
}

// GetMe returns the current authenticated user
func (c *Client) GetMe() (*User, error) {
	var user User
	if err := c.get("/api/users/me", &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetProjects returns all projects
func (c *Client) GetProjects() ([]Project, error) {
	var resp APIResponse
	if err := c.get("/api/projects", &resp); err != nil {
		return nil, err
	}
	return extractItems[Project](resp)
}

// GetProject returns a project by ID
func (c *Client) GetProject(projectID string) (*Project, error) {
	var resp struct {
		Item     Project                `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/projects/%s", projectID), &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(req CreateProjectRequest) (*Project, error) {
	var resp struct {
		Item Project `json:"item"`
	}
	if err := c.post("/api/projects", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	return c.delete(fmt.Sprintf("/api/projects/%s", projectID))
}

// GetBoards returns all boards for a project
// Note: Boards are included in the project response, so we get the project and extract boards from included
func (c *Client) GetBoards(projectID string) ([]Board, error) {
	var resp struct {
		Item     Project                `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/projects/%s", projectID), &resp); err != nil {
		return nil, err
	}
	
	// Extract boards from included
	if boardsData, ok := resp.Included["boards"]; ok {
		boardsJSON, err := json.Marshal(boardsData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal boards: %w", err)
		}
		
		var boards []Board
		if err := json.Unmarshal(boardsJSON, &boards); err != nil {
			return nil, fmt.Errorf("failed to unmarshal boards: %w", err)
		}
		return boards, nil
	}
	
	return []Board{}, nil
}

// GetBoard returns a board by ID
func (c *Client) GetBoard(boardID string) (*Board, error) {
	var resp struct {
		Item     Board                  `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/boards/%s", boardID), &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// CreateBoard creates a new board
// Note: Boards are created via /api/projects/{projectId}/boards endpoint and require a position
func (c *Client) CreateBoard(req CreateBoardRequest) (*Board, error) {
	var resp struct {
		Item Board `json:"item"`
	}
	// Position is required - use default if not provided
	position := req.Position
	if position == 0 {
		position = 65535 // Default position
	}
	
	// Create request body without projectId (it's in the URL)
	requestBody := map[string]interface{}{
		"name":     req.Name,
		"position": position,
	}
	if req.Description != "" {
		requestBody["description"] = req.Description
	}
	if err := c.post(fmt.Sprintf("/api/projects/%s/boards", req.ProjectID), requestBody, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// DeleteBoard deletes a board
func (c *Client) DeleteBoard(boardID string) error {
	return c.delete(fmt.Sprintf("/api/boards/%s", boardID))
}

// GetLists returns all lists for a board
// Note: Lists are included in the board response, so we get the board and extract lists from included
func (c *Client) GetLists(boardID string) ([]List, error) {
	var resp struct {
		Item     Board                  `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/boards/%s", boardID), &resp); err != nil {
		return nil, err
	}
	
	// Extract lists from included
	if listsData, ok := resp.Included["lists"]; ok {
		listsJSON, err := json.Marshal(listsData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal lists: %w", err)
		}
		
		var lists []List
		if err := json.Unmarshal(listsJSON, &lists); err != nil {
			return nil, fmt.Errorf("failed to unmarshal lists: %w", err)
		}
		return lists, nil
	}
	
	return []List{}, nil
}

// GetList returns a list by ID
func (c *Client) GetList(listID string) (*List, error) {
	var resp struct {
		Item     List                   `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/lists/%s", listID), &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// CreateList creates a new list
// Note: Lists are created via /api/boards/{boardId}/lists endpoint and require a position
func (c *Client) CreateList(req CreateListRequest) (*List, error) {
	// Position is required - use default if not provided
	position := req.Position
	if position == 0 {
		position = 65535 // Default position
	}
	
	// Create request body without boardId (it's in the URL)
	requestBody := map[string]interface{}{
		"name":     req.Name,
		"position": position,
	}
	
	var resp struct {
		Item List `json:"item"`
	}
	if err := c.post(fmt.Sprintf("/api/boards/%s/lists", req.BoardID), requestBody, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// DeleteList deletes a list
func (c *Client) DeleteList(listID string) error {
	return c.delete(fmt.Sprintf("/api/lists/%s", listID))
}

// GetCards returns all cards for a list
// Note: Cards are included in the board response. We need to find which board contains this list.
// Since we can't reliably get the list directly, we'll need the boardId. 
// For now, we'll get all boards and search for the one containing this list, then get its cards.
// Alternatively, if boardId is known, use GetBoards and filter.
func (c *Client) GetCards(listID string) ([]Card, error) {
	// Try to get the list first - if it works, use the boardId from it
	var listResp struct {
		Item     List                   `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	
	// Try getting list - if it fails with HTML, we'll need another approach
	err := c.get(fmt.Sprintf("/api/lists/%s", listID), &listResp)
	var boardID string
	
	if err != nil {
		// List endpoint returned HTML, so we need to find the board another way
		// Get all projects and search through boards to find the one with this list
		projects, err := c.GetProjects()
		if err != nil {
			return nil, fmt.Errorf("failed to get projects to find board: %w", err)
		}
		
		// Search through projects and boards to find the list
		for _, project := range projects {
			boards, err := c.GetBoards(project.ID)
			if err != nil {
				continue
			}
			for _, board := range boards {
				lists, err := c.GetLists(board.ID)
				if err != nil {
					continue
				}
				for _, list := range lists {
					if list.ID == listID {
						boardID = board.ID
						break
					}
				}
				if boardID != "" {
					break
				}
			}
			if boardID != "" {
				break
			}
		}
		
		if boardID == "" {
			return []Card{}, nil
		}
	} else {
		boardID = listResp.Item.BoardID
		if boardID == "" {
			return []Card{}, nil
		}
	}
	
	// Get the board which includes all cards
	var boardResp struct {
		Item     Board                  `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/boards/%s", boardID), &boardResp); err != nil {
		return nil, fmt.Errorf("failed to get board %s: %w", boardID, err)
	}
	
	// Extract cards from included and filter by listId
	if cardsData, ok := boardResp.Included["cards"]; ok {
		cardsJSON, err := json.Marshal(cardsData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal cards: %w", err)
		}
		
		var allCards []Card
		if err := json.Unmarshal(cardsJSON, &allCards); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cards: %w", err)
		}
		
		// Filter cards by listId
		var filteredCards []Card
		for _, card := range allCards {
			if card.ListID == listID {
				filteredCards = append(filteredCards, card)
			}
		}
		
		return filteredCards, nil
	}
	
	return []Card{}, nil
}

// GetCard returns a card by ID
func (c *Client) GetCard(cardID string) (*Card, error) {
	var resp struct {
		Item     Card                   `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/cards/%s", cardID), &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// CreateCard creates a new card
// Note: Cards are created via /api/lists/{listId}/cards endpoint
func (c *Client) CreateCard(req CreateCardRequest) (*Card, error) {
	var resp struct {
		Item Card `json:"item"`
	}
	// Position is required - use default if not provided
	position := req.Position
	if position == 0 {
		position = 65535 // Default position
	}
	
	// Create request body without listId (it's in the URL)
	requestBody := map[string]interface{}{
		"name":     req.Name,
		"position": position,
	}
	if req.Description != "" {
		requestBody["description"] = req.Description
	}
	if req.DueDate != nil {
		requestBody["dueDate"] = req.DueDate.Format(time.RFC3339)
	}
	if err := c.post(fmt.Sprintf("/api/lists/%s/cards", req.ListID), requestBody, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// UpdateCard updates a card
func (c *Client) UpdateCard(cardID string, req UpdateCardRequest) (*Card, error) {
	var resp struct {
		Item Card `json:"item"`
	}
	if err := c.patch(fmt.Sprintf("/api/cards/%s", cardID), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// DeleteCard deletes a card
func (c *Client) DeleteCard(cardID string) error {
	return c.delete(fmt.Sprintf("/api/cards/%s", cardID))
}

// MoveCard moves a card to a different list
func (c *Client) MoveCard(cardID, listID string, position float64) (*Card, error) {
	req := UpdateCardRequest{
		ListID:   &listID,
		Position: &position,
	}
	return c.UpdateCard(cardID, req)
}

// GetTasks returns all tasks for a card
// Note: Tasks are included in the card response
func (c *Client) GetTasks(cardID string) ([]Task, error) {
	var resp struct {
		Item     Card                   `json:"item"`
		Included map[string]interface{} `json:"included,omitempty"`
	}
	if err := c.get(fmt.Sprintf("/api/cards/%s", cardID), &resp); err != nil {
		return nil, err
	}
	
	// Extract tasks from included
	if tasksData, ok := resp.Included["tasks"]; ok {
		tasksJSON, err := json.Marshal(tasksData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tasks: %w", err)
		}
		
		var tasks []Task
		if err := json.Unmarshal(tasksJSON, &tasks); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tasks: %w", err)
		}
		return tasks, nil
	}
	
	return []Task{}, nil
}

// CreateTask creates a new task
// Note: Tasks are created via /api/cards/{cardId}/tasks endpoint
func (c *Client) CreateTask(req CreateTaskRequest) (*Task, error) {
	var resp struct {
		Item Task `json:"item"`
	}
	// Position is required - use default if not provided
	position := req.Position
	if position == 0 {
		position = 65535 // Default position
	}
	
	// Create request body without cardId (it's in the URL)
	requestBody := map[string]interface{}{
		"name":     req.Name,
		"position": position,
	}
	if err := c.post(fmt.Sprintf("/api/cards/%s/tasks", req.CardID), requestBody, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// UpdateTask updates a task
func (c *Client) UpdateTask(taskID string, req UpdateTaskRequest) (*Task, error) {
	var resp struct {
		Item Task `json:"item"`
	}
	if err := c.patch(fmt.Sprintf("/api/tasks/%s", taskID), req, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(taskID string) error {
	return c.delete(fmt.Sprintf("/api/tasks/%s", taskID))
}

// GetComments returns all comments for a card
// Note: Comments endpoint may return HTML, so we try the endpoint first, and if it fails,
// we check if comments are in the card's included section
func (c *Client) GetComments(cardID string) ([]Comment, error) {
	// Try the comments endpoint first
	var resp APIResponse
	err := c.get(fmt.Sprintf("/api/cards/%s/comments", cardID), &resp)
	
	if err != nil {
		// Endpoint returned HTML, try getting from card's included section
		var cardResp struct {
			Item     Card                   `json:"item"`
			Included map[string]interface{} `json:"included,omitempty"`
		}
		if err := c.get(fmt.Sprintf("/api/cards/%s", cardID), &cardResp); err != nil {
			return nil, fmt.Errorf("failed to get card: %w", err)
		}
		
		// Extract comments from included
		if commentsData, ok := cardResp.Included["comments"]; ok {
			commentsJSON, err := json.Marshal(commentsData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal comments: %w", err)
			}
			
			var comments []Comment
			if err := json.Unmarshal(commentsJSON, &comments); err != nil {
				return nil, fmt.Errorf("failed to unmarshal comments: %w", err)
			}
			return comments, nil
		}
		
		return []Comment{}, nil
	}
	
	return extractItems[Comment](resp)
}

// CreateComment creates a new comment
func (c *Client) CreateComment(req CreateCommentRequest) (*Comment, error) {
	var resp struct {
		Item Comment `json:"item"`
	}
	if err := c.post("/api/comments", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// DeleteComment deletes a comment
func (c *Client) DeleteComment(commentID string) error {
	return c.delete(fmt.Sprintf("/api/comments/%s", commentID))
}

// GetStopwatch returns the stopwatch for a card
func (c *Client) GetStopwatch(cardID string) (*Stopwatch, error) {
	var resp struct {
		Item Stopwatch `json:"item"`
	}
	if err := c.get(fmt.Sprintf("/api/cards/%s/stopwatch", cardID), &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// StartStopwatch starts the stopwatch for a card
func (c *Client) StartStopwatch(cardID string) (*Stopwatch, error) {
	var resp struct {
		Item Stopwatch `json:"item"`
	}
	if err := c.post(fmt.Sprintf("/api/cards/%s/stopwatch/start", cardID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// StopStopwatch stops the stopwatch for a card
func (c *Client) StopStopwatch(cardID string) (*Stopwatch, error) {
	var resp struct {
		Item Stopwatch `json:"item"`
	}
	if err := c.post(fmt.Sprintf("/api/cards/%s/stopwatch/stop", cardID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

// ResetStopwatch resets the stopwatch for a card
func (c *Client) ResetStopwatch(cardID string) (*Stopwatch, error) {
	var resp struct {
		Item Stopwatch `json:"item"`
	}
	if err := c.post(fmt.Sprintf("/api/cards/%s/stopwatch/reset", cardID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Item, nil
}

