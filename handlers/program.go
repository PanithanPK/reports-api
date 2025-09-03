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

// @Summary List programs
// @Description Get list of all programs with pagination
// @Tags programs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/list [get]
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
		SELECT sp.id, sp.name, IFNULL(sp.type, 0) as type_id, IFNULL(it.name, '') as type_name, sp.created_at, sp.updated_at, sp.deleted_at, sp.created_by, sp.updated_by, sp.deleted_by 
		FROM systems_program sp
		LEFT JOIN issue_types it ON sp.type = it.id
		WHERE deleted_at IS NULL
		ORDER BY sp.id DESC
		LIMIT ? OFFSET ?
	`, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query programs"})
	}
	defer rows.Close()

	var programs []models.Program
	for rows.Next() {
		var p models.Program
		err := rows.Scan(&p.ID, &p.Name, &p.TypeID, &p.TypeName, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy)
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

// @Summary Create program
// @Description Create a new program
// @Tags programs
// @Accept json
// @Produce json
// @Param program body models.ProgramRequest true "Program data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/create [post]
func CreateProgramHandler(c *fiber.Ctx) error {
	var req models.ProgramRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Priority == nil || *req.Priority == 0 {
		priority := 2
		req.Priority = &priority
	}

	res, err := db.DB.Exec(`INSERT INTO systems_program (name, priority, type, created_by) VALUES (?, ?, ?, ?)`, req.Name, req.Priority, req.TypeID, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert program"})
	}

	id, _ := res.LastInsertId()
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// @Summary Update program
// @Description Update an existing program
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Program ID"
// @Param program body models.ProgramRequest true "Program data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/update/{id} [put]
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

	_, err = db.DB.Exec(`UPDATE systems_program SET name=?, type=?, updated_by=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Name, req.TypeID, req.UpdatedBy, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update program"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// @Summary Delete program
// @Description Delete a program
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Program ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/delete/{id} [delete]
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

// @Summary Search programs
// @Description Search programs by query string with pagination
// @Tags programs
// @Accept json
// @Produce json
// @Param query path string true "Search query (use 'all' for all programs)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/list/{query} [get]
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
		SELECT sp.id, sp.name, IFNULL(sp.type, 0), IFNULL(it.name, ''), sp.created_at, sp.updated_at, sp.deleted_at, sp.created_by, sp.updated_by, sp.deleted_by
		FROM systems_program sp
		LEFT JOIN issue_types it ON sp.type = it.id
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
		err := rows.Scan(&p.ID, &p.Name, &p.TypeID, &p.TypeName, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy)
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

// @Summary Get program details
// @Description Get detailed information about a specific program
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Program ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/program/{id} [get]
func GetProgramDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var program models.Program
	err = db.DB.QueryRow(`
		SELECT sp.id, sp.name, IFNULL(sp.type, 0), IFNULL(it.name, ''), sp.created_at, sp.updated_at, sp.deleted_at, sp.created_by, sp.updated_by, sp.deleted_by
		FROM systems_program sp
		LEFT JOIN issue_types it ON sp.type = it.id
		WHERE sp.id=? AND sp.deleted_at IS NULL
	`, id).Scan(
		&program.ID, &program.Name, &program.TypeID, &program.TypeName, &program.CreatedAt, &program.UpdatedAt,
		&program.DeletedAt, &program.CreatedBy, &program.UpdatedBy, &program.DeletedBy,
	)

	if err != nil {
		log.Printf("Error fetching program details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Program not found"})
	}

	log.Printf("Getting program details Success for ID: %d", id)
	return c.JSON(fiber.Map{"success": true, "data": program})
}

// @Summary Get program types
// @Description Get all program types
// @Tags programs
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/type/list [get]
func GETTypeProgramHandler(c *fiber.Ctx) error {
	rows, err := db.DB.Query(`SELECT id, name FROM issue_types ORDER BY id`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query program types"})
	}
	defer rows.Close()

	var types []fiber.Map
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Printf("Error scanning program type: %v", err)
			continue
		}
		types = append(types, fiber.Map{"id": id, "name": name})
	}
	log.Printf("Getting issue types Success")
	return c.JSON(fiber.Map{"success": true, "data": types})
}

// @Summary Add program type
// @Description Add a new program type
// @Tags programs
// @Accept json
// @Produce json
// @Param type body models.Type true "Type data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/type/create [post]
func AddTypeProgramHandler(c *fiber.Ctx) error {
	var req models.Type
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := db.DB.Exec(`INSERT INTO issue_types (name) VALUES (?)`, req.Name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert program type"})
	}

	id, _ := res.LastInsertId()
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// @Summary Update program type
// @Description Update an existing program type
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Type ID"
// @Param type body models.Type true "Type data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/type/update/{id} [post]
func UpdateTypeProgramHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req models.Type
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err = db.DB.Exec(`UPDATE issue_types SET name=? WHERE id=?`, req.Name, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update program type"})
	}

	return c.JSON(fiber.Map{"success": true})
}

// @Summary Delete program type
// @Description Delete a program type
// @Tags programs
// @Accept json
// @Produce json
// @Param id path string true "Type ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/type/delete/{id} [delete]
func DeleteTypeHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = db.DB.Exec(`DELETE FROM issue_types WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete program type"})
	}

	log.Printf("Deleted program type ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

// @Summary Search program types
// @Description Search program types by query string
// @Tags programs
// @Accept json
// @Produce json
// @Param query path string true "Search query (use 'all' for all types)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/program/type/list/{query} [get]
func GetTypeWithQueryHandler(c *fiber.Ctx) error {
	query := c.Params("query")

	// If query is empty or "all", return all types
	if query == "" || query == "all" {
		return GETTypeProgramHandler(c)
	}

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

	// Get search results
	rows, err := db.DB.Query(`
		SELECT id, name FROM issue_types
		WHERE name LIKE ?
		ORDER BY id
	`, searchPattern)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search types"})
	}
	defer rows.Close()

	var types []fiber.Map
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Printf("Error scanning type: %v", err)
			continue
		}
		types = append(types, fiber.Map{"id": id, "name": name})
	}

	log.Printf("Searching types with query: %s, found %d results", query, len(types))
	return c.JSON(fiber.Map{"success": true, "data": types})
}
