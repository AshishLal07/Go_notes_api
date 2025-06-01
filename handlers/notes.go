package handlers

import (
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"notes-api/config"
	"notes-api/middleware"
	"notes-api/models"
	"notes-api/utils"
)

// NotesHandler handles note-related operations
type NotesHandler struct{}

// NewNotesHandler creates a new notes handler
func NewNotesHandler() *NotesHandler {
	return &NotesHandler{}
}

// CreateNote creates a new note for the authenticated user
func (h *NotesHandler) CreateNote(c *fiber.Ctx) error {
	var req models.NoteCreateRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Create new note
	note := models.Note{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	// Save note to database
	if err := config.GetDB().Create(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create note",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Note created successfully",
		"data":    note.ToResponse(),
	})
}

// GetNotes retrieves all notes for the authenticated user with pagination and search
func (h *NotesHandler) GetNotes(c *fiber.Ctx) error {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Parse search parameter
	search := strings.TrimSpace(c.Query("search", ""))

	// Build query
	query := config.GetDB().Where("user_id = ?", userID)

	// Add search filter if provided
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	var total int64
	if err := query.Model(&models.Note{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to count notes",
		})
	}

	// Calculate pagination
	offset := (page - 1) * perPage
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	// Fetch notes with pagination
	var notes []models.Note
	if err := query.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&notes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch notes",
		})
	}

	// Convert to response format
	var noteResponses []models.NoteResponse
	for _, note := range notes {
		noteResponses = append(noteResponses, note.ToResponse())
	}

	// Build paginated response
	response := models.PaginatedNotesResponse{
		Notes:       noteResponses,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Notes retrieved successfully",
		"data":    response,
	})
}

// GetNote retrieves a specific note by ID (only if it belongs to the authenticated user)
func (h *NotesHandler) GetNote(c *fiber.Ctx) error {
	// Get note ID from URL parameter
	noteID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid note ID",
		})
	}

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Find note
	var note models.Note
	if err := config.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Note not found",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Note retrieved successfully",
		"data":    note.ToResponse(),
	})
}

// UpdateNote updates a specific note (only if it belongs to the authenticated user)
func (h *NotesHandler) UpdateNote(c *fiber.Ctx) error {
	// Get note ID from URL parameter
	noteID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid note ID",
		})
	}

	var req models.NoteUpdateRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Find note
	var note models.Note
	if err := config.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Note not found",
		})
	}

	// Update note
	note.Title = req.Title
	note.Content = req.Content

	if err := config.GetDB().Save(&note).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update note",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Note updated successfully",
		"data":    note.ToResponse(),
	})
}

// DeleteNote deletes a specific note (only if it belongs to the authenticated user)
func (h *NotesHandler) DeleteNote(c *fiber.Ctx) error {
	// Get note ID from URL parameter
	noteID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid note ID",
		})
	}

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Find and delete note
	result := config.GetDB().Where("id = ? AND user_id = ?", noteID, userID).Delete(&models.Note{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete note",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Note not found",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Note deleted successfully",
	})
}
