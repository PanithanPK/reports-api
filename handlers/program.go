package handlers

import (
	"log"
	"net/url"
	"reports-api/db"
	"reports-api/models"
	"reports-api/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ListProgramsHandler returns a handler for listing all programs with pagination
func ListProgramsHandler(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// Get total count
	var total int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM systems_program WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count programs"})
	}

	// Get paginated data
	rows, err := db.DB.Query(`
		SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by 
		FROM systems_program 
		WHERE deleted_at IS NULL 
		LIMIT ? OFFSET ?
	`, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query programs"})
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, p)
	}

	log.Printf("Getting programs Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    programs,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}

// CreateProgramHandler returns a handler for creating a new program
func CreateProgramHandler(c *fiber.Ctx) error {
	var req models.ProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := db.DB.Exec(`INSERT INTO systems_program (name, created_by) VALUES (?, ?)`, req.Name, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert program"})
	}

	id, _ := res.LastInsertId()
	log.Printf("Inserted new program: %s", req.Name)
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// UpdateProgramHandler returns a handler for updating an existing program
func UpdateProgramHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req models.ProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err = db.DB.Exec(`UPDATE systems_program SET name=?, updated_by=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Name, req.UpdatedBy, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update program"})
	}

	log.Printf("Updating program ID: %d with name: %s", id, req.Name)
	return c.JSON(fiber.Map{"success": true})
}

// DeleteProgramHandler returns a handler for deleting a program (soft delete)
func DeleteProgramHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = db.DB.Exec(`DELETE FROM systems_program WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete program"})
	}

	log.Printf("Deleted program ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

func ListProgramsQueryHandler(c *fiber.Ctx) error {
	query := c.Params("query")

	// If query is empty or "all", return all programs with pagination
	if query == "" || query == "all" {
		return ListProgramsHandler(c)
	}

	// Otherwise, search programs
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// URL decode for Thai language support
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		decodedQuery = query
	}

	// Clean query
	decodedQuery = strings.TrimSpace(decodedQuery)
	decodedQuery = strings.ReplaceAll(decodedQuery, "  ", " ")
	decodedQuery = strings.ReplaceAll(decodedQuery, "%", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "_", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "'", "")
	decodedQuery = strings.ReplaceAll(decodedQuery, "\"", "")

	searchPattern := "%" + decodedQuery + "%"

	// Get total count
	var total int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM systems_program
		WHERE deleted_at IS NULL AND name LIKE ?
	`, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results with pagination
	rows, err := db.DB.Query(`
		SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
		FROM systems_program
		WHERE deleted_at IS NULL AND name LIKE ?
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, searchPattern, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search programs"})
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, p)
	}

	log.Printf("Searching programs with query: %s, found %d results", query, len(programs))
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    programs,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}
