package planka

import "time"

// User represents a Planka user
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// Project represents a Planka project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Boards      []Board   `json:"boards,omitempty"`
}

// Board represents a Planka board
type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProjectID   string    `json:"projectId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Lists       []List    `json:"lists,omitempty"`
}

// List represents a Planka list
type List struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BoardID   string    `json:"boardId"`
	Position  float64   `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Cards     []Card    `json:"cards,omitempty"`
}

// Card represents a Planka card
type Card struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ListID      string    `json:"listId"`
	Position    float64   `json:"position"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tasks       []Task    `json:"tasks,omitempty"`
	Comments    []Comment `json:"comments,omitempty"`
	Labels      []Label   `json:"labels,omitempty"`
}

// Task represents a task within a card
type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CardID    string    `json:"cardId"`
	Position  float64   `json:"position"`
	IsCompleted bool    `json:"isCompleted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Comment represents a comment on a card
type Comment struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CardID    string    `json:"cardId"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Label represents a label on a card
type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Stopwatch represents a time tracking stopwatch
type Stopwatch struct {
	ID        string     `json:"id"`
	CardID    string     `json:"cardId"`
	StartedAt *time.Time `json:"startedAt,omitempty"`
	Duration  int64      `json:"duration"` // Duration in seconds
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// CreateBoardRequest represents a request to create a board
type CreateBoardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ProjectID   string `json:"projectId"`
}

// CreateListRequest represents a request to create a list
type CreateListRequest struct {
	Name     string  `json:"name"`
	BoardID  string  `json:"boardId"`
	Position float64 `json:"position"` // Position is required by the API
}

// CreateCardRequest represents a request to create a card
type CreateCardRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ListID      string     `json:"listId"`
	Position    float64    `json:"position,omitempty"`
	DueDate     *time.Time `json:"dueDate,omitempty"`
}

// UpdateCardRequest represents a request to update a card
type UpdateCardRequest struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	ListID      *string     `json:"listId,omitempty"`
	Position    *float64    `json:"position,omitempty"`
	DueDate     *time.Time  `json:"dueDate,omitempty"`
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	Name     string  `json:"name"`
	CardID   string  `json:"cardId"`
	Position float64 `json:"position,omitempty"`
}

// UpdateTaskRequest represents a request to update a task
type UpdateTaskRequest struct {
	Name        *string `json:"name,omitempty"`
	IsCompleted *bool   `json:"isCompleted,omitempty"`
	Position    *float64 `json:"position,omitempty"`
}

// CreateCommentRequest represents a request to create a comment
type CreateCommentRequest struct {
	Text   string `json:"text"`
	CardID string `json:"cardId"`
}

