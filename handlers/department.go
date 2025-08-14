package handlers

import (
	"log"
	"net/url"
	"reports-api/db"
	"reports-api/models"
	"reports-api/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ListDepartmentsHandler returns a handler for listing all departments with pagination
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

// CreateDepartmentHandler returns a handler for creating a new department
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

// UpdateDepartmentHandler returns a handler for updating an existing department
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

// DeleteDepartmentHandler returns a handler for deleting a department
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

// GetDepartmentDetailHandler returns detailed information about a specific department
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

// SearchDepartmentsHandler returns a handler for searching departments
func SearchDepartmentsHandler(c *fiber.Ctx) error {
	query := c.Params("query")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Query parameter is required"})
	}

	// URL decode for Thai language support
	decodedQuery, err := url.QueryUnescape(query)
	if err != nil {
		decodedQuery = query
	}
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

	// Get search results
	rows, err := db.DB.Query(`
		SELECT d.id, d.name, d.branch_id, IFNULL(b.name, '') as branch_name, d.created_at, d.updated_at, d.deleted_at
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		WHERE d.deleted_at IS NULL
		AND (d.name LIKE ? OR b.name LIKE ?)
		ORDER BY d.id DESC
	`, searchPattern, searchPattern)
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

	log.Printf("Searching departments with query: %s", query)
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    departments,
	})
}
