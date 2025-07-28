package handlers

import (
	"log"
	"reports-api/db"
	"reports-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// ListDepartmentsHandler returns a handler for listing all departments
func ListDepartmentsHandler(c *fiber.Ctx) error {
	rows, err := db.DB.Query(`
   SELECT d.id, d.name, d.branch_id, b.name as branch_name, d.created_at, d.updated_at, d.deleted_at
   FROM departments d
   LEFT JOIN branches b ON d.branch_id = b.id
   WHERE d.deleted_at IS NULL
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
	return c.JSON(fiber.Map{"success": true, "data": departments})
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
		SELECT d.id, d.name, d.branch_id, b.name, d.created_at, d.updated_at 
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