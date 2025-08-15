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

// ListBranchesHandler returns a handler for listing all branches with pagination
func ListBranchesHandler(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// Get total count
	var total int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM branches WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count branches"})
	}

	// Get paginated data
	rows, err := db.DB.Query(`
		SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by 
		FROM branches 
		WHERE deleted_at IS NULL 
		ORDER BY id DESC 
		LIMIT ? OFFSET ?
	`, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query branches"})
	}
	defer rows.Close()

	var branches []models.Branch
	for rows.Next() {
		var b models.Branch
		err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy, &b.DeletedBy)
		if err != nil {
			log.Printf("Error scanning branch: %v", err)
			continue
		}
		branches = append(branches, b)
	}

	log.Printf("Getting branches Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    branches,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}

// CreateBranchHandler returns a handler for creating a new branch
func CreateBranchHandler(c *fiber.Ctx) error {
	var req models.BranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := db.DB.Exec(`INSERT INTO branches (name, created_by, updated_by) VALUES (?, ?, ?)`, req.Name, req.CreatedBy, req.UpdatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert branch"})
	}

	id, _ := res.LastInsertId()
	log.Printf("Inserted new branch: %s", req.Name)
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// UpdateBranchHandler returns a handler for updating an existing branch
func UpdateBranchHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req models.BranchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err = db.DB.Exec(`UPDATE branches SET name=?, updated_by=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Name, req.UpdatedBy, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update branch"})
	}

	log.Printf("Updating branch ID: %d with name: %s", id, req.Name)
	return c.JSON(fiber.Map{"success": true})
}

// DeleteBranchHandler returns a handler for deleting a branch
func DeleteBranchHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = db.DB.Exec(`DELETE FROM branches WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete branch"})
	}

	log.Printf("Deleted branch ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

// GetBranchDetailHandler returns detailed information about a specific branch
func GetBranchDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var branchDetail models.BranchDetail
	err = db.DB.QueryRow(`
		SELECT id, name, created_at, updated_at 
		FROM branches 
		WHERE id = ? AND deleted_at IS NULL
	`, id).Scan(&branchDetail.ID, &branchDetail.Name, &branchDetail.CreatedAt, &branchDetail.UpdatedAt)

	if err != nil {
		log.Printf("Error fetching branch details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "Branch not found"})
	}

	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM departments 
		WHERE branch_id = ? AND deleted_at IS NULL
	`, id).Scan(&branchDetail.DepartmentsCount)
	if err != nil {
		log.Printf("Error counting departments: %v", err)
	}

	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM ip_phones ip
		JOIN departments d ON ip.department_id = d.id
		WHERE d.branch_id = ? AND ip.deleted_at IS NULL
	`, id).Scan(&branchDetail.IPPhonesCount)

	if err != nil {
		log.Printf("Error counting IP phones: %v", err)
	}

	log.Printf("Getting branch details Success for ID: %d", id)
	return c.JSON(fiber.Map{"success": true, "data": branchDetail})
}

func ListBranchesQueryHandler(c *fiber.Ctx) error {
	query := c.Params("query")

	// If query is empty or "all", return all branches with pagination
	if query == "" || query == "all" {
		return ListBranchesHandler(c)
	}

	// Otherwise, search branches
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
		SELECT COUNT(*) FROM branches
		WHERE deleted_at IS NULL AND name LIKE ?
	`, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results with pagination
	rows, err := db.DB.Query(`
		SELECT id, name, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
		FROM branches
		WHERE deleted_at IS NULL AND name LIKE ?
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`, searchPattern, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search branches"})
	}
	defer rows.Close()

	var branches []models.Branch
	for rows.Next() {
		var b models.Branch
		err := rows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt, &b.CreatedBy, &b.UpdatedBy, &b.DeletedBy)
		if err != nil {
			log.Printf("Error scanning branch: %v", err)
			continue
		}
		branches = append(branches, b)
	}

	log.Printf("Searching branches with query: %s, found %d results", query, len(branches))
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    branches,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}
