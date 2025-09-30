package handlers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"reports-api/db"
	"time"

	"github.com/gofiber/fiber/v2"
)

func IpphonesExportCsv(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filename := fmt.Sprintf("IPPhones_%s.csv", time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	c.Response().BodyWriter().Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	headers := []string{"ID", "Number", "Name", "Department", "Branch", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to write headers"})
	}

	query := `
		SELECT 
			ip.id, 
			ip.number, 
			ip.name, 
			d.name as department_name,
			b.name as branch_name,
			DATE_ADD(ip.created_at, INTERVAL 7 HOUR) as created_at,
			DATE_ADD(ip.updated_at, INTERVAL 7 HOUR) as updated_at
		FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE ip.deleted_at IS NULL`

	var args []interface{}
	// แก้ไข: ถ้ามี startDate หรือ endDate อย่างใดอย่างหนึ่ง ให้ใส่เงื่อนไขวันที่เสมอ
	if startDate != "" || endDate != "" {
		if startDate != "" && endDate != "" {
			query += ` AND DATE(ip.created_at) BETWEEN ? AND ?`
			args = append(args, startDate, endDate)
		} else if startDate != "" {
			query += ` AND DATE(ip.created_at) >= ?`
			args = append(args, startDate)
		} else if endDate != "" {
			query += ` AND DATE(ip.created_at) <= ?`
			args = append(args, endDate)
		}
	}
	query += ` ORDER BY ip.id`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query data"})
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var number, name, departmentName, branchName, createdAt, updatedAt sql.NullString

		err := rows.Scan(&id, &number, &name, &departmentName, &branchName, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}

		record := []string{
			fmt.Sprintf("%d", id),
			number.String,
			name.String,
			departmentName.String,
			branchName.String,
			createdAt.String,
			updatedAt.String,
		}

		if err := writer.Write(record); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to write record"})
		}
	}

	return nil
}

func DepartmentsExportCsv(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filename := fmt.Sprintf("Departments_%s.csv", time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	c.Response().BodyWriter().Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	headers := []string{"ID", "Name", "Branch", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to write headers"})
	}

	query := `
		SELECT 
			d.id, 
			d.name, 
			b.name as branch_name,
			DATE_ADD(d.created_at, INTERVAL 7 HOUR) as created_at,
			DATE_ADD(d.updated_at, INTERVAL 7 HOUR) as updated_at
		FROM departments d
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE d.deleted_at IS NULL`

	var args []interface{}
	// แก้ไข: ถ้ามี startDate หรือ endDate อย่างใดอย่างหนึ่ง ให้ใส่เงื่อนไขวันที่เสมอ
	if startDate != "" || endDate != "" {
		if startDate != "" && endDate != "" {
			query += ` AND DATE(d.created_at) BETWEEN ? AND ?`
			args = append(args, startDate, endDate)
		} else if startDate != "" {
			query += ` AND DATE(d.created_at) >= ?`
			args = append(args, startDate)
		} else if endDate != "" {
			query += ` AND DATE(d.created_at) <= ?`
			args = append(args, endDate)
		}
	}
	query += ` ORDER BY d.id`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query data"})
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, branchName, createdAt, updatedAt sql.NullString

		err := rows.Scan(&id, &name, &branchName, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}

		record := []string{
			fmt.Sprintf("%d", id),
			name.String,
			branchName.String,
			createdAt.String,
			updatedAt.String,
		}

		if err := writer.Write(record); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to write record"})
		}
	}

	return nil
}

func BranchExportCsv(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filename := fmt.Sprintf("Branches_%s.csv", time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	c.Response().BodyWriter().Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	headers := []string{"ID", "Name", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to write headers"})
	}

	query := `
		SELECT 
			id, 
			name,
			DATE_ADD(created_at, INTERVAL 7 HOUR) as created_at,
			DATE_ADD(updated_at, INTERVAL 7 HOUR) as updated_at
		FROM branches
		WHERE deleted_at IS NULL`

	var args []interface{}
	// แก้ไข: ถ้ามี startDate หรือ endDate อย่างใดอย่างหนึ่ง ให้ใส่เงื่อนไขวันที่เสมอ
	if startDate != "" || endDate != "" {
		if startDate != "" && endDate != "" {
			query += ` AND DATE(created_at) BETWEEN ? AND ?`
			args = append(args, startDate, endDate)
		} else if startDate != "" {
			query += ` AND DATE(created_at) >= ?`
			args = append(args, startDate)
		} else if endDate != "" {
			query += ` AND DATE(created_at) <= ?`
			args = append(args, endDate)
		}
	}
	query += ` ORDER BY id`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query data"})
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, createdAt, updatedAt sql.NullString

		err := rows.Scan(&id, &name, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}

		record := []string{
			fmt.Sprintf("%d", id),
			name.String,
			createdAt.String,
			updatedAt.String,
		}

		if err := writer.Write(record); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to write record"})
		}
	}

	return nil
}

func SystemExportCsv(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filename := fmt.Sprintf("Systems_%s.csv", time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	c.Response().BodyWriter().Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	headers := []string{"ID", "Name", "Priority", "Issue Type", "Created At", "Updated At"}
	if err := writer.Write(headers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to write headers"})
	}

	query := `
		SELECT 
			sp.id, 
			sp.name, 
			sp.priority,
			it.name as issue_type_name,
			DATE_ADD(sp.created_at, INTERVAL 7 HOUR) as created_at,
			DATE_ADD(sp.updated_at, INTERVAL 7 HOUR) as updated_at
		FROM systems_program sp
		LEFT JOIN issue_types it ON sp.type = it.id
		WHERE sp.deleted_at IS NULL`

	var args []interface{}
	// แก้ไข: ถ้ามี startDate หรือ endDate อย่างใดอย่างหนึ่ง ให้ใส่เงื่อนไขวันที่เสมอ
	if startDate != "" || endDate != "" {
		if startDate != "" && endDate != "" {
			query += ` AND DATE(sp.created_at) BETWEEN ? AND ?`
			args = append(args, startDate, endDate)
		} else if startDate != "" {
			query += ` AND DATE(sp.created_at) >= ?`
			args = append(args, startDate)
		} else if endDate != "" {
			query += ` AND DATE(sp.created_at) <= ?`
			args = append(args, endDate)
		}
	}
	query += ` ORDER BY sp.id`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query data"})
	}
	defer rows.Close()

	for rows.Next() {
		var id, priority int
		var name, issueTypeName, createdAt, updatedAt sql.NullString

		err := rows.Scan(&id, &name, &priority, &issueTypeName, &createdAt, &updatedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}

		record := []string{
			fmt.Sprintf("%d", id),
			name.String,
			fmt.Sprintf("%d", priority),
			issueTypeName.String,
			createdAt.String,
			updatedAt.String,
		}

		if err := writer.Write(record); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to write record"})
		}
	}

	return nil
}

func TasksExportCsv(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filename := fmt.Sprintf("Tasks_%s.csv", time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	c.Response().BodyWriter().Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	headers := []string{"ID", "Ticket No", "Phone Name", "Issue Type", "System/Issue", "Branch", "Department", "Description", "Reported By", "Assigned To", "Solution", "Status", "Progress Notes", "Created At", "Updated At", "Resolved At"}
	if err := writer.Write(headers); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to write headers"})
	}

	query := `
		SELECT 
			t.id,
			t.ticket_no,
			COALESCE(ip.name, CONCAT('Phone #', ip.number)) as phone_name,
			it.name as issue_type_name,
			CASE 
				WHEN t.system_id = 0 OR t.system_id IS NULL THEN t.issue_else
				ELSE sp.name 
			END as system_name,
			b.name as branch_name,
			d.name as department_name,
			t.text,
			t.reported_by,
			r.name as assignto_name,
			res.text as solution_text,
			CASE 
				WHEN t.status = 0 THEN 'รอดำเนินการ'
				WHEN t.status = 1 THEN 'กำลังดำเนินการ'
				WHEN t.status = 2 THEN 'เสร็จสิ้นแล้ว'
				ELSE 'ไม่ระบุ'
			END as status_text,
			GROUP_CONCAT(DISTINCT p.progress_text ORDER BY p.created_at SEPARATOR ' , ') as progress_notes,
			DATE_ADD(t.created_at, INTERVAL 7 HOUR) as created_at,
			DATE_ADD(t.updated_at, INTERVAL 7 HOUR) as updated_at,
			CASE WHEN t.resolved_at IS NOT NULL THEN DATE_ADD(t.resolved_at, INTERVAL 7 HOUR) ELSE NULL END as resolved_at
		FROM tasks t
		LEFT JOIN ip_phones ip ON t.phone_id = ip.id
		LEFT JOIN issue_types it ON t.issue_type = it.id
		LEFT JOIN systems_program sp ON t.system_id = sp.id AND t.system_id != 0
		LEFT JOIN departments d ON t.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		LEFT JOIN responsibilities r ON t.assignto_id = r.id
		LEFT JOIN resolutions res ON t.solution_id = res.id
		LEFT JOIN progress p ON t.id = p.task_id
		WHERE t.deleted_at IS NULL`

	var args []interface{}
	// แก้ไข: ถ้ามี startDate หรือ endDate อย่างใดอย่างหนึ่ง ให้ใส่เงื่อนไขวันที่เสมอ
	if startDate != "" || endDate != "" {
		if startDate != "" && endDate != "" {
			query += ` AND DATE(t.created_at) BETWEEN ? AND ?`
			args = append(args, startDate, endDate)
		} else if startDate != "" {
			query += ` AND DATE(t.created_at) >= ?`
			args = append(args, startDate)
		} else if endDate != "" {
			query += ` AND DATE(t.created_at) <= ?`
			args = append(args, endDate)
		}
	}
	query += ` GROUP BY t.id ORDER BY t.id`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query data"})
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var ticketNo, phoneName, issueTypeName, systemName, branchName, departmentName, text, reportedBy, assigntoName, solutionText, statusText, progressNotes, createdAt, updatedAt, resolvedAt sql.NullString

		err := rows.Scan(&id, &ticketNo, &phoneName, &issueTypeName, &systemName, &branchName, &departmentName, &text, &reportedBy, &assigntoName, &solutionText, &statusText, &progressNotes, &createdAt, &updatedAt, &resolvedAt)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to scan data"})
		}

		record := []string{
			fmt.Sprintf("%d", id),
			ticketNo.String,
			phoneName.String,
			issueTypeName.String,
			systemName.String,
			branchName.String,
			departmentName.String,
			text.String,
			reportedBy.String,
			assigntoName.String,
			solutionText.String,
			statusText.String,
			progressNotes.String,
			createdAt.String,
			updatedAt.String,
			resolvedAt.String,
		}

		if err := writer.Write(record); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to write record"})
		}
	}

	return nil
}
