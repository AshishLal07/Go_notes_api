package models

import (
	"time"

	"gorm.io/gorm"
)

// Note represents a note in the system
type Note struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null;size:200" validate:"required,min=1,max=200"`
	Content   string         `json:"content" gorm:"type:text" validate:"required,min=1"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// NoteCreateRequest represents the note creation request payload
type NoteCreateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// NoteUpdateRequest represents the note update request payload
type NoteUpdateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// NoteResponse represents the note response
type NoteResponse struct {
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	UserID    uint         `json:"user_id"`
	User      UserResponse `json:"user,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ToResponse converts Note to NoteResponse
func (n *Note) ToResponse() NoteResponse {
	response := NoteResponse{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		UserID:    n.UserID,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
	
	if n.User.ID != 0 {
		response.User = n.User.ToResponse()
	}
	
	return response
}

// PaginatedNotesResponse represents paginated notes response
type PaginatedNotesResponse struct {
	Notes       []NoteResponse `json:"notes"`
	Total       int64          `json:"total"`
	Page        int            `json:"page"`
	PerPage     int            `json:"per_page"`
	TotalPages  int            `json:"total_pages"`
	HasNext     bool           `json:"has_next"`
	HasPrevious bool           `json:"has_previous"`
}