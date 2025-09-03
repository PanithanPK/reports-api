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

// @Summary List departments
// @Description Get list of all departments with pagination
// @Tags departments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/department/list [get]
func ListDepartmentsHandler(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// Get total count
	var total int
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM departments d
		WHERE d.deleted_at IS NULL
	`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count departments"})
	}

	// Get paginated data
	rows, err := db.DB.Query(`
		SELECT d.id, d.name, d.branch_id, b.name as branch_name, d.created_at, d.updated_at, d.deleted_at
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE d.deleted_at IS NULL
		ORDER BY d.id DESC
		LIMIT ? OFFSET ?
	`, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query departments"})
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.BranchID, &d.BranchName, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
		if err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		departments = append(departments, d)
	}

	log.Printf("Getting departments Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    departments,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}

// @Summary Create department
// @Description Create a new department
// @Tags departments
// @Accept json
// @Produce json
// @Param department body models.DepartmentRequest true "Department data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/department/create [post]
func CreateDepartmentHandler(c *fiber.Ctx) error {
	var req models.DepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := db.DB.Exec(`INSERT INTO departments (name, branch_id) VALUES (?, ?)`, req.Name, req.BranchID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert department"})
	}

	id, _ := res.LastInsertId()
	log.Printf("Inserted new department: %s", req.Name)
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// @Summary Update department
// @Description Update an existing department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param department body models.DepartmentRequest true "Department data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/department/update/{id} [put]
func UpdateDepartmentHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req models.DepartmentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err = db.DB.Exec(`UPDATE departments SET name=?, branch_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Name, req.BranchID, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update department"})
	}

	log.Printf("Updating department ID: %d with name: %s", id, req.Name)
	return c.JSON(fiber.Map{"success": true})
}

// @Summary Delete department
// @Description Delete a department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/department/delete/{id} [delete]
func DeleteDepartmentHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = db.DB.Exec(`DELETE FROM departments WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete department"})
	}

	log.Printf("Deleted department ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

// @Summary Get department details
// @Description Get detailed information about a specific department
// @Tags departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/department/{id} [get]
func GetDepartmentDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var departmentDetail models.DepartmentDetail
	err = db.DB.QueryRow(`
		SELECT d.id, d.name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name, d.created_at, d.updated_at 
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE d.id = ? AND d.deleted_at IS NULL
	`, id).Scan(
		&departmentDetail.ID,
		&departmentDetail.Name,
		&departmentDetail.BranchID,
		&departmentDetail.BranchName,
		&departmentDetail.CreatedAt,
		&departmentDetail.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error fetching department details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Department not found"})
	}

	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM ip_phones 
		WHERE department_id = ? AND deleted_at IS NULL
	`, id).Scan(&departmentDetail.IPPhonesCount)

	if err != nil {
		log.Printf("Error counting IP phones: %v", err)
	}

	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM tasks t
		JOIN ip_phones ip ON t.phone_id = ip.id
		WHERE ip.department_id = ? AND t.deleted_at IS NULL
	`, id).Scan(&departmentDetail.TasksCount)

	if err != nil {
		log.Printf("Error counting tasks: %v", err)
	}

	log.Printf("Getting department details Success for ID: %d", id)
	return c.JSON(fiber.Map{"success": true, "data": departmentDetail})
}

// @Summary Get all departments
// @Description Get all departments without pagination
// @Tags departments
// @Accept json
// @Produce json
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/department/listall [get]
func AllDepartmentsHandler(c *fiber.Ctx) error {
	// Get total count
	var total int
	err := db.DB.QueryRow(`
		SELECT COUNT(*) FROM departments d
		WHERE d.deleted_at IS NULL
	`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count departments"})
	}

	// Get paginated data
	rows, err := db.DB.Query(`
		SELECT d.id, d.name, d.branch_id, b.name as branch_name, d.created_at, d.updated_at, d.deleted_at
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE d.deleted_at IS NULL
		ORDER BY d.id DESC
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query departments"})
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.BranchID, &d.BranchName, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
		if err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		departments = append(departments, d)
	}

	log.Printf("Getting departments Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    departments,
	})
}

// @Summary Search departments
// @Description Search departments by query string with pagination
// @Tags departments
// @Accept json
// @Produce json
// @Param query path string true "Search query (use 'all' for all departments)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/department/list/{query} [get]
func ListDepartmentsQueryHandler(c *fiber.Ctx) error {
	query := c.Params("query")

	// If query is empty or "all", return all departments with pagination
	if query == "" || query == "all" {
		return ListDepartmentsHandler(c)
	}

	// Otherwise, search departments
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
		SELECT COUNT(*) FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		WHERE d.deleted_at IS NULL
		AND (d.name LIKE ? OR b.name LIKE ?)
	`, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results with pagination
	rows, err := db.DB.Query(`
		SELECT d.id, d.name, d.branch_id, IFNULL(b.name, '') as branch_name, d.created_at, d.updated_at, d.deleted_at
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		WHERE d.deleted_at IS NULL
		AND (d.name LIKE ? OR b.name LIKE ?)
		ORDER BY d.id DESC
		LIMIT ? OFFSET ?
	`, searchPattern, searchPattern, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search departments"})
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.BranchID, &d.BranchName, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
		if err != nil {
			log.Printf("Error scanning department: %v", err)
			continue
		}
		departments = append(departments, d)
	}

	log.Printf("Searching departments with query: %s, found %d results", query, len(departments))
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    departments,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}
