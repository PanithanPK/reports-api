package handlers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"reports-api/db"
	"reports-api/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ExportToCSV(c *fiber.Ctx) error {
	tableName := c.Query("table", "tasks")

	filename := fmt.Sprintf("Report_%s_%s.csv", tableName, time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	// เขียน BOM สำหรับ UTF-8 เพื่อให้ Excel รองรับภาษาไทย
	c.Response().BodyWriter().Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	err := exportTableToCSV(writer, db.DB, tableName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export table " + tableName + ": " + err.Error()})
	}

	return nil
}

func IpphonesExportCsv(c *fiber.Ctx) error {
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
		WHERE ip.deleted_at IS NULL
		ORDER BY ip.id`

	rows, err := db.DB.Query(query)
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
		WHERE d.deleted_at IS NULL
		ORDER BY d.id`

	rows, err := db.DB.Query(query)
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
		WHERE deleted_at IS NULL
		ORDER BY id`

	rows, err := db.DB.Query(query)
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
		WHERE sp.deleted_at IS NULL
		ORDER BY sp.id`

	rows, err := db.DB.Query(query)
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

func exportTableToCSV(writer *csv.Writer, db *sql.DB, tableName string) error {
	var query string
	var headers []string

	if tableName == "tasks" {
		query = `
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
		WHERE t.deleted_at IS NULL
		GROUP BY t.id`

		headers = []string{
			"ID", "Ticket No", "Phone Name", "Issue Type", "System/Issue",
			"Branch", "Department", "Description", "Reported By", "Assigned To",
			"Solution", "Status", "Progress Notes", "Created At", "Updated At", "Resolved At",
		}
	} else {
		query, headers = buildFilteredQuery(tableName)
	}

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write data rows
	for rows.Next() {
		if tableName == "tasks" {
			var task models.DataTask
			var createdAtStr, updatedAtStr string
			var resolvedAtStr, progressNotesStr sql.NullString
			err := rows.Scan(&task.ID, &task.TicketNo, &task.PhoneName, &task.IssueTypeName,
				&task.SystemName, &task.BranchName, &task.DepartmentName, &task.Text,
				&task.ReportedBy, &task.AssigntoName, &task.SolutionText, &task.StatusText,
				&progressNotesStr, &createdAtStr, &updatedAtStr, &resolvedAtStr)
			if err != nil {
				return err
			}

			var ticketNo, phoneName, issueTypeName, systemName, branchName, departmentName, text, reportedBy, assigntoName, solutionText, statusText string
			if task.TicketNo != nil {
				ticketNo = *task.TicketNo
			}
			if task.PhoneName != nil {
				phoneName = *task.PhoneName
			}
			if task.IssueTypeName != nil {
				issueTypeName = *task.IssueTypeName
			}
			if task.SystemName != nil {
				systemName = *task.SystemName
			}
			if task.BranchName != nil {
				branchName = *task.BranchName
			}
			if task.DepartmentName != nil {
				departmentName = *task.DepartmentName
			}
			if task.Text != nil {
				text = *task.Text
			}
			if task.ReportedBy != nil {
				reportedBy = *task.ReportedBy
			}
			if task.AssigntoName != nil {
				assigntoName = *task.AssigntoName
			}
			if task.SolutionText != nil {
				solutionText = *task.SolutionText
			}
			if task.StatusText != nil {
				statusText = *task.StatusText
			}

			record := []string{
				fmt.Sprintf("%d", task.ID), ticketNo, phoneName, issueTypeName,
				systemName, branchName, departmentName, text,
				reportedBy, assigntoName, solutionText, statusText,
				progressNotesStr.String, createdAtStr, updatedAtStr, resolvedAtStr.String,
			}

			if err := writer.Write(record); err != nil {
				return err
			}
		} else {
			var record []string
			switch tableName {
			case "issue_types":
				var item models.DataIssueType
				var createdAtStr string
				err := rows.Scan(&item.ID, &item.Name, &createdAtStr)
				if err != nil {
					return err
				}
				record = []string{fmt.Sprintf("%d", item.ID), item.Name, createdAtStr}
			case "systems_program":
				var item models.DataSystemProgram
				var createdAtStr, updatedAtStr string
				var typeName sql.NullString
				err := rows.Scan(&item.ID, &item.Name, &item.Priority, &typeName, &createdAtStr, &updatedAtStr)
				if err != nil {
					return err
				}
				var name string
				var priority int
				if item.Name != nil {
					name = *item.Name
				}
				if item.Priority != nil {
					priority = *item.Priority
				}
				record = []string{fmt.Sprintf("%d", item.ID), name, fmt.Sprintf("%d", priority), typeName.String, createdAtStr, updatedAtStr}
			case "branches":
				var item models.DataBranch
				var createdAtStr, updatedAtStr string
				err := rows.Scan(&item.ID, &item.Name, &createdAtStr, &updatedAtStr)
				if err != nil {
					return err
				}
				var name string
				if item.Name != nil {
					name = *item.Name
				}
				record = []string{fmt.Sprintf("%d", item.ID), name, createdAtStr, updatedAtStr}
			case "departments":
				var item models.DataDepartment
				var createdAtStr, updatedAtStr string
				var branchName sql.NullString
				err := rows.Scan(&item.ID, &item.Name, &branchName, &createdAtStr, &updatedAtStr)
				if err != nil {
					return err
				}
				var name string
				if item.Name != nil {
					name = *item.Name
				}
				record = []string{fmt.Sprintf("%d", item.ID), name, branchName.String, createdAtStr, updatedAtStr}
			case "ip_phones":
				var item models.DataIPPhone
				var createdAtStr, updatedAtStr string
				var departmentName, branchName sql.NullString
				err := rows.Scan(&item.ID, &item.Number, &item.Name, &departmentName, &branchName, &createdAtStr, &updatedAtStr)
				if err != nil {
					return err
				}
				var number, name string
				if item.Number != nil {
					number = *item.Number
				}
				if item.Name != nil {
					name = *item.Name
				}
				record = []string{fmt.Sprintf("%d", item.ID), number, name, departmentName.String, branchName.String, createdAtStr, updatedAtStr}
			case "resolutions":
				var item models.DataResolution
				var resolvedAtStr sql.NullString
				var updatedAtStr string
				var ticketNo sql.NullString
				err := rows.Scan(&item.ID, &ticketNo, &item.Text, &resolvedAtStr, &updatedAtStr)
				if err != nil {
					return err
				}
				var text string
				if item.Text != nil {
					text = *item.Text
				}
				record = []string{fmt.Sprintf("%d", item.ID), ticketNo.String, text, resolvedAtStr.String, updatedAtStr}
			case "responsibilities":
				var item models.DataResponsibility
				var createdAtStr, updatedAtStr string
				err := rows.Scan(&item.ID, &item.TelegramUsername, &item.Name, &createdAtStr, &updatedAtStr)
				if err != nil {
					return err
				}
				var telegramUsername, name string
				if item.TelegramUsername != nil {
					telegramUsername = *item.TelegramUsername
				}
				if item.Name != nil {
					name = *item.Name
				}
				record = []string{fmt.Sprintf("%d", item.ID), telegramUsername, name, createdAtStr, updatedAtStr}
			default:
				values := make([]interface{}, len(headers))
				valuePtrs := make([]interface{}, len(headers))
				for i := range values {
					valuePtrs[i] = &values[i]
				}
				if err := rows.Scan(valuePtrs...); err != nil {
					return err
				}
				record = make([]string, len(values))
				for i, val := range values {
					if val != nil {
						record[i] = fmt.Sprintf("%v", val)
					}
				}
			}

			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildFilteredQuery(tableName string) (string, []string) {
	switch tableName {
	case "issue_types":
		return `SELECT id, name, DATE_ADD(created_at, INTERVAL 7 HOUR) as created_at FROM issue_types`,
			[]string{"ID", "Name", "Created At"}

	case "systems_program":
		return `SELECT sp.id, sp.name, sp.priority, it.name as type_name, DATE_ADD(sp.created_at, INTERVAL 7 HOUR) as created_at, DATE_ADD(sp.updated_at, INTERVAL 7 HOUR) as updated_at FROM systems_program sp LEFT JOIN issue_types it ON sp.type = it.id WHERE sp.deleted_at IS NULL`,
			[]string{"ID", "Name", "Priority", "Type", "Created At", "Updated At"}

	case "branches":
		return `SELECT id, name, DATE_ADD(created_at, INTERVAL 7 HOUR) as created_at, DATE_ADD(updated_at, INTERVAL 7 HOUR) as updated_at FROM branches WHERE deleted_at IS NULL`,
			[]string{"ID", "Name", "Created At", "Updated At"}

	case "departments":
		return `SELECT d.id, d.name, b.name as branch_name, DATE_ADD(d.created_at, INTERVAL 7 HOUR) as created_at, DATE_ADD(d.updated_at, INTERVAL 7 HOUR) as updated_at FROM departments d LEFT JOIN branches b ON d.branch_id = b.id WHERE d.deleted_at IS NULL`,
			[]string{"ID", "Name", "Branch", "Created At", "Updated At"}

	case "ip_phones":
		return `SELECT ip.id, ip.number, ip.name, d.name as department_name, b.name as branch_name, DATE_ADD(ip.created_at, INTERVAL 7 HOUR) as created_at, DATE_ADD(ip.updated_at, INTERVAL 7 HOUR) as updated_at FROM ip_phones ip LEFT JOIN departments d ON ip.department_id = d.id LEFT JOIN branches b ON d.branch_id = b.id WHERE ip.deleted_at IS NULL`,
			[]string{"ID", "Number", "Name", "Department", "Branch", "Created At", "Updated At"}

	case "resolutions":
		return `SELECT r.id, t.ticket_no, r.text, DATE_ADD(r.resolved_at, INTERVAL 7 HOUR) as resolved_at, DATE_ADD(r.updated_at, INTERVAL 7 HOUR) as updated_at FROM resolutions r LEFT JOIN tasks t ON r.tasks_id = t.id WHERE r.deleted_at IS NULL`,
			[]string{"ID", "Ticket No", "Text", "Resolved At", "Updated At"}

	case "responsibilities":
		return `SELECT id, telegram_username, name, DATE_ADD(created_at, INTERVAL 7 HOUR) as created_at, DATE_ADD(updated_at, INTERVAL 7 HOUR) as updated_at FROM responsibilities`,
			[]string{"ID", "Telegram Username", "Name", "Created At", "Updated At"}

	default:
		return fmt.Sprintf("SELECT * FROM %s", tableName), []string{}
	}
}
