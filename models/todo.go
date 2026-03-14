package models

import "time"

// Todo represents a todo item in the database
// @Description Todo item model
type Todo struct {
	ID        uint      `gorm:"primaryKey" json:"id" example:"1"`
	Title     string    `gorm:"not null" json:"title" binding:"required" example:"Buy milk"`
	Priority  int       `gorm:"default:0" json:"priority" example:"1" enums:"0,1,2,3"`
	Completed bool      `gorm:"default:false" json:"completed" example:"false"`
	DueDate   time.Time `json:"due_date" example:"2023-01-01T00:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateTodoInput represents the input for creating a todo item
// @Description Input for creating a todo item
type CreateTodoInput struct {
	Title    string    `json:"title" binding:"required" example:"Buy milk"`
	Priority int       `json:"priority" example:"1" enums:"0,1,2,3"`
	DueDate  time.Time `json:"due_date" binding:"required" example:"2023-01-01T00:00:00Z"`
}

// UpdateTodoInput represents the input for updating a todo item
// @Description Input for updating a todo item
type UpdateTodoInput struct {
	Title     string     `json:"title" example:"Buy milk"`
	Priority  int        `json:"priority" example:"1" enums:"0,1,2,3"`
	Completed bool       `json:"completed" example:"true"`
	DueDate   *time.Time `json:"due_date" example:"2023-01-01T00:00:00Z"`
}

// APIError represents an error response from the API
// @Description API error response
type APIError struct {
	Error string `json:"error" example:"Internal server error"`
}
