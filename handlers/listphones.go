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

// ListIPPhonesHandler returns a handler for listing all IP phones with pagination
func ListIPPhonesHandler(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)
	offset := utils.CalculateOffset(pagination.Page, pagination.Limit)

	// Get total count
	var total int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM ip_phones WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count ip_phones"})
	}

	// Get paginated data
	rows, err := db.DB.Query(`
		SELECT ip.id, ip.number, ip.name, IFNULL(ip.department_id,0) as department_id,
			   IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
			   ip.created_at, ip.updated_at, ip.deleted_at, ip.created_by, ip.updated_by, ip.deleted_by
		FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE ip.deleted_at IS NULL
		ORDER BY ip.id DESC
		LIMIT ? OFFSET ?
	`, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query ip_phones"})
	}
	defer rows.Close()

	var phones []models.IPPhone
	for rows.Next() {
		var p models.IPPhone
		err := rows.Scan(
			&p.ID, &p.Number, &p.Name, &p.DepartmentID,
			&p.DepartmentName, &p.BranchID, &p.BranchName,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy,
		)
		if err != nil {
			log.Printf("Error scanning ip_phone: %v", err)
			continue
		}
		phones = append(phones, p)
	}

	log.Printf("Getting IP phones Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    phones,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}

// CreateIPPhoneHandler returns a handler for creating a new IP phone
func CreateIPPhoneHandler(c *fiber.Ctx) error {
	var req models.IPPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	res, err := db.DB.Exec(`INSERT INTO ip_phones (number, name, department_id, created_by) VALUES (?, ?, ?, ?)`, req.Number, req.Name, req.DepartmentID, req.CreatedBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to insert ip_phone"})
	}

	log.Printf("Inserted new IP phone: %d", req.Number)
	id, _ := res.LastInsertId()
	return c.JSON(fiber.Map{"success": true, "id": id})
}

// UpdateIPPhoneHandler returns a handler for updating an existing IP phone
func UpdateIPPhoneHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var req models.IPPhoneRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	_, err = db.DB.Exec(`UPDATE ip_phones SET number=?, name=?, department_id=?, updated_by=?, updated_at=CURRENT_TIMESTAMP WHERE id=? AND deleted_at IS NULL`, req.Number, req.Name, req.DepartmentID, req.UpdatedBy, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update ip_phone"})
	}

	log.Printf("Updating IP phone ID: %d with number: %d", id, req.Number)
	return c.JSON(fiber.Map{"success": true})
}

// DeleteIPPhoneHandler returns a handler for deleting an IP phone
func DeleteIPPhoneHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	_, err = db.DB.Exec(`DELETE FROM ip_phones WHERE id=?`, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete ip_phone"})
	}

	log.Printf("Deleted IP phone ID: %d", id)
	return c.JSON(fiber.Map{"success": true})
}

func AllIPPhonesHandler(c *fiber.Ctx) error {

	// Get total count
	var total int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM ip_phones WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count ip_phones"})
	}

	// Get paginated data
	rows, err := db.DB.Query(`
		SELECT ip.id, ip.number, ip.name, IFNULL(ip.department_id,0) as department_id,
			   IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
			   ip.created_at, ip.updated_at, ip.deleted_at, ip.created_by, ip.updated_by, ip.deleted_by
		FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE ip.deleted_at IS NULL
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query ip_phones"})
	}
	defer rows.Close()

	var phones []models.IPPhone
	for rows.Next() {
		var p models.IPPhone
		err := rows.Scan(
			&p.ID, &p.Number, &p.Name, &p.DepartmentID,
			&p.DepartmentName, &p.BranchID, &p.BranchName,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy,
		)
		if err != nil {
			log.Printf("Error scanning ip_phone: %v", err)
			continue
		}
		phones = append(phones, p)
	}

	log.Printf("Getting IP phones Success")
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    phones,
	})
}

func ListIPPhonesQueryHandler(c *fiber.Ctx) error {
	query := c.Params("query")

	// If query is empty or "all", return all IP phones with pagination
	if query == "" || query == "all" {
		return ListIPPhonesHandler(c)
	}

	// Otherwise, search IP phones
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
		SELECT COUNT(*) FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id AND d.deleted_at IS NULL
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		WHERE ip.deleted_at IS NULL
		AND (ip.number LIKE ? OR ip.name LIKE ? OR d.name LIKE ? OR b.name LIKE ?)
	`, searchPattern, searchPattern, searchPattern, searchPattern).Scan(&total)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count search results"})
	}

	// Get search results with pagination
	rows, err := db.DB.Query(`
		SELECT ip.id, ip.number, ip.name, IFNULL(ip.department_id,0) as department_id,
			   IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
			   ip.created_at, ip.updated_at, ip.deleted_at, ip.created_by, ip.updated_by, ip.deleted_by
		FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id AND d.deleted_at IS NULL
		LEFT JOIN branches b ON d.branch_id = b.id AND b.deleted_at IS NULL
		WHERE ip.deleted_at IS NULL
		AND (ip.number LIKE ? OR ip.name LIKE ? OR d.name LIKE ? OR b.name LIKE ?)
		ORDER BY ip.id DESC
		LIMIT ? OFFSET ?
	`, searchPattern, searchPattern, searchPattern, searchPattern, pagination.Limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to search ip_phones"})
	}
	defer rows.Close()

	var phones []models.IPPhone
	for rows.Next() {
		var p models.IPPhone
		err := rows.Scan(
			&p.ID, &p.Number, &p.Name, &p.DepartmentID,
			&p.DepartmentName, &p.BranchID, &p.BranchName,
			&p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.CreatedBy, &p.UpdatedBy, &p.DeletedBy,
		)
		if err != nil {
			log.Printf("Error scanning ip_phone: %v", err)
			continue
		}
		phones = append(phones, p)
	}

	log.Printf("Searching IP phones with query: %s, found %d results", query, len(phones))
	return c.JSON(models.PaginatedResponse{
		Success: true,
		Data:    phones,
		Pagination: models.PaginationResponse{
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			Total:      total,
			TotalPages: utils.CalculateTotalPages(total, pagination.Limit),
		},
	})
}

func GetIPPhonesDetailHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid id"})
	}

	var ipPhone models.IPPhone
	err = db.DB.QueryRow(`
		SELECT ip.id, ip.number, ip.name, IFNULL(ip.department_id,0) as department_id,
			   IFNULL(d.name, '') as department_name, IFNULL(d.branch_id, 0) as branch_id, IFNULL(b.name, '') as branch_name,
			   ip.created_at, ip.updated_at, ip.deleted_at, ip.created_by, ip.updated_by, ip.deleted_by
		FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE ip.id = ? AND ip.deleted_at IS NULL
	`, id).Scan(
		&ipPhone.ID, &ipPhone.Number, &ipPhone.Name, &ipPhone.DepartmentID,
		&ipPhone.DepartmentName, &ipPhone.BranchID, &ipPhone.BranchName,
		&ipPhone.CreatedAt, &ipPhone.UpdatedAt, &ipPhone.DeletedAt,
		&ipPhone.CreatedBy, &ipPhone.UpdatedBy, &ipPhone.DeletedBy,
	)

	if err != nil {
		log.Printf("Error fetching IP phone details: %v", err)
		return c.Status(404).JSON(fiber.Map{"error": "IP phone not found"})
	}

	log.Printf("Getting IP phone details Success for ID: %d", id)
	return c.JSON(fiber.Map{"success": true, "data": ipPhone})
}
