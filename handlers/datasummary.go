package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"reports-api/db"
	"reports-api/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func ExportToExcel(c *fiber.Ctx) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing Excel file:", err)
		}
	}()

	// สร้าง styles ที่จะใช้
	err := createExcelStyles(f)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create styles: " + err.Error()})
	}

	tables := []string{"tasks", "issue_types", "systems_program", "branches", "departments", "ip_phones", "resolutions", "responsibilities"}

	for i, table := range tables {
		if i == 0 {
			f.SetSheetName("Sheet1", table)
		} else {
			_, err := f.NewSheet(table)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create sheet: " + err.Error()})
			}
		}

		err := exportTableToSheet(f, db.DB, table)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to export table " + table + ": " + err.Error()})
		}
	}

	filename := fmt.Sprintf("Report_db_datasummary_%s.xlsx", time.Now().Add(7*time.Hour).Format("20060102_150405"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename="+filename)

	if err := f.Write(c.Response().BodyWriter()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to write Excel file: " + err.Error()})
	}

	return nil
}

func createExcelStyles(f *excelize.File) error {
	// 1. Header Style (สีน้ำเงินเข้ม - Professional Blue)
	_, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"1F4E79"}, // Dark Blue
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 2},
			{Type: "right", Color: "FFFFFF", Style: 2},
			{Type: "top", Color: "FFFFFF", Style: 2},
			{Type: "bottom", Color: "FFFFFF", Style: 2},
		},
	})
	if err != nil {
		return err
	}

	// 2. Data Style แถวคี่ (Light Blue)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Color:  "333333",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"E7F3FF"}, // Light Blue
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "B8CCE4", Style: 1},
			{Type: "right", Color: "B8CCE4", Style: 1},
			{Type: "top", Color: "B8CCE4", Style: 1},
			{Type: "bottom", Color: "B8CCE4", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 3. Data Style แถวคู่ (White)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   10,
			Color:  "333333",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFFFF"}, // White
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "B8CCE4", Style: 1},
			{Type: "right", Color: "B8CCE4", Style: 1},
			{Type: "top", Color: "B8CCE4", Style: 1},
			{Type: "bottom", Color: "B8CCE4", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 4. Status Style - เสร็จสิ้น (Green Success)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"28A745"}, // Success Green
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "1E7E34", Style: 1},
			{Type: "right", Color: "1E7E34", Style: 1},
			{Type: "top", Color: "1E7E34", Style: 1},
			{Type: "bottom", Color: "1E7E34", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 5. Status Style - รอดำเนินการ (Orange Warning)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FD7E14"}, // Orange Warning
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "E55100", Style: 1},
			{Type: "right", Color: "E55100", Style: 1},
			{Type: "top", Color: "E55100", Style: 1},
			{Type: "bottom", Color: "E55100", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 6. Priority High Style (Red)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"DC3545"}, // Danger Red
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "BD2130", Style: 1},
			{Type: "right", Color: "BD2130", Style: 1},
			{Type: "top", Color: "BD2130", Style: 1},
			{Type: "bottom", Color: "BD2130", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 7. Priority Medium Style (Yellow)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Color:  "333333",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFC107"}, // Warning Yellow
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "E0A800", Style: 1},
			{Type: "right", Color: "E0A800", Style: 1},
			{Type: "top", Color: "E0A800", Style: 1},
			{Type: "bottom", Color: "E0A800", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 8. Priority Low Style (Light Green)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Color:  "FFFFFF",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"6F42C1"}, // Purple Info
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "5A32A3", Style: 1},
			{Type: "right", Color: "5A32A3", Style: 1},
			{Type: "top", Color: "5A32A3", Style: 1},
			{Type: "bottom", Color: "5A32A3", Style: 1},
		},
	})
	if err != nil {
		return err
	}

	// 9. Date/Time Style (Light Gray)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:   9,
			Color:  "495057",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"F8F9FA"}, // Light Gray
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "DEE2E6", Style: 1},
			{Type: "right", Color: "DEE2E6", Style: 1},
			{Type: "top", Color: "DEE2E6", Style: 1},
			{Type: "bottom", Color: "DEE2E6", Style: 1},
		},
		NumFmt: 22, // Date format
	})
	if err != nil {
		return err
	}

	// 10. Number/ID Style (Light Blue Background)
	_, err = f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   10,
			Color:  "0D47A1",
			Family: "Calibri",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"E3F2FD"}, // Light Blue Background
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "90CAF9", Style: 1},
			{Type: "right", Color: "90CAF9", Style: 1},
			{Type: "top", Color: "90CAF9", Style: 1},
			{Type: "bottom", Color: "90CAF9", Style: 1},
		},
	})

	return err
}

func exportTableToSheet(f *excelize.File, db *sql.DB, tableName string) error {
	var query string
	var columns []string

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
				WHEN t.status = 1 THEN 'เสร็จสิ้นแล้ว'
				ELSE 'ไม่ระบุ'
			END as status_text,
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
		WHERE t.deleted_at IS NULL`

		columns = []string{
			"ID", "Ticket No", "Phone Name", "Issue Type", "System/Issue",
			"Branch", "Department", "Description", "Reported By", "Assigned To",
			"Solution", "Status", "Created At", "Updated At", "Resolved At",
		}
	} else {
		query, columns = buildFilteredQuery(tableName)
	}

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// ตั้งค่า Header
	for i, col := range columns {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(tableName, cell, col)
	}

	// Apply Header Style
	headerRange := fmt.Sprintf("A1:%s1", string(rune('A'+len(columns)-1)))
	f.SetCellStyle(tableName, "A1", headerRange, 1) // Style ID 1 = Header Style

	// Set row height for header
	f.SetRowHeight(tableName, 1, 30)

	rowNum := 2
	for rows.Next() {
		if tableName == "tasks" {
			var task models.DataTask
			var createdAtStr, updatedAtStr string
			var resolvedAtStr sql.NullString
			err := rows.Scan(&task.ID, &task.TicketNo, &task.PhoneName, &task.IssueTypeName,
				&task.SystemName, &task.BranchName, &task.DepartmentName, &task.Text,
				&task.ReportedBy, &task.AssigntoName, &task.SolutionText, &task.StatusText,
				&createdAtStr, &updatedAtStr, &resolvedAtStr)
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

			values := []interface{}{task.ID, ticketNo, phoneName, issueTypeName,
				systemName, branchName, departmentName, text,
				reportedBy, assigntoName, solutionText, statusText,
				createdAtStr, updatedAtStr, resolvedAtStr.String}

			for i, val := range values {
				cell := fmt.Sprintf("%s%d", string(rune('A'+i)), rowNum)
				if val != nil {
					f.SetCellValue(tableName, cell, val)
				}

				// Apply styling based on column type and content
				switch i {
				case 0: // ID Column
					f.SetCellStyle(tableName, cell, cell, 10) // Number/ID style
				case 11: // Status column
					statusVal := fmt.Sprintf("%v", val)
					if strings.Contains(statusVal, "เสร็จสิ้นแล้ว") {
						f.SetCellStyle(tableName, cell, cell, 4) // Green style
					} else if strings.Contains(statusVal, "รอดำเนินการ") {
						f.SetCellStyle(tableName, cell, cell, 5) // Orange style
					} else {
						// Apply alternating row colors for unknown status
						if rowNum%2 == 0 {
							f.SetCellStyle(tableName, cell, cell, 3) // White
						} else {
							f.SetCellStyle(tableName, cell, cell, 2) // Light Blue
						}
					}
				case 12, 13, 14: // Date columns
					f.SetCellStyle(tableName, cell, cell, 9) // Date style
				default:
					// Apply alternating row colors
					if rowNum%2 == 0 {
						f.SetCellStyle(tableName, cell, cell, 3) // White
					} else {
						f.SetCellStyle(tableName, cell, cell, 2) // Light Blue
					}
				}
			}
		} else {
			var values []interface{}
			switch tableName {
			case "issue_types":
				var item models.DataIssueType
				var createdAtStr string
				err := rows.Scan(&item.ID, &item.Name, &createdAtStr)
				if err != nil {
					return err
				}
				values = []interface{}{item.ID, item.Name, createdAtStr}
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
				values = []interface{}{item.ID, name, priority, typeName.String, createdAtStr, updatedAtStr}
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
				values = []interface{}{item.ID, name, createdAtStr, updatedAtStr}
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
				values = []interface{}{item.ID, name, branchName.String, createdAtStr, updatedAtStr}
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
				values = []interface{}{item.ID, number, name, departmentName.String, branchName.String, createdAtStr, updatedAtStr}
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
				values = []interface{}{item.ID, ticketNo.String, text, resolvedAtStr.String, updatedAtStr}
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
				values = []interface{}{item.ID, telegramUsername, name, createdAtStr, updatedAtStr}
			default:
				values = make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}
				if err := rows.Scan(valuePtrs...); err != nil {
					return err
				}
			}

			// Apply styling for non-task tables
			for i, val := range values {
				cell := fmt.Sprintf("%s%d", string(rune('A'+i)), rowNum)
				if val != nil {
					f.SetCellValue(tableName, cell, val)
				}

				// Apply specific styling based on column content
				colName := strings.ToLower(columns[i])
				switch {
				case i == 0: // ID column
					f.SetCellStyle(tableName, cell, cell, 10) // Number/ID style
				case strings.Contains(colName, "priority"):
					// Apply priority colors for systems_program table
					if tableName == "systems_program" {
						priorityVal := fmt.Sprintf("%v", val)
						switch priorityVal {
						case "1", "高": // High priority
							f.SetCellStyle(tableName, cell, cell, 6) // Red
						case "2", "中": // Medium priority
							f.SetCellStyle(tableName, cell, cell, 7) // Yellow
						case "3", "低": // Low priority
							f.SetCellStyle(tableName, cell, cell, 8) // Purple
						default:
							// Apply normal alternating colors
							if rowNum%2 == 0 {
								f.SetCellStyle(tableName, cell, cell, 3) // White
							} else {
								f.SetCellStyle(tableName, cell, cell, 2) // Light Blue
							}
						}
					}
				case strings.Contains(colName, "created") || strings.Contains(colName, "updated") || strings.Contains(colName, "resolved"):
					f.SetCellStyle(tableName, cell, cell, 9) // Date style
				default:
					// Apply alternating row colors
					if rowNum%2 == 0 {
						f.SetCellStyle(tableName, cell, cell, 3) // White
					} else {
						f.SetCellStyle(tableName, cell, cell, 2) // Light Blue
					}
				}
			}
		}

		// Set row height
		f.SetRowHeight(tableName, rowNum, 20)
		rowNum++
	}

	// Auto-adjust column widths และ freeze panes
	for i, col := range columns {
		colLetter := string(rune('A' + i))

		// กำหนดความกว้างคอลัมน์ตามประเภทข้อมูล
		var width float64
		switch {
		case strings.Contains(strings.ToLower(col), "id"):
			width = 10
		case strings.Contains(strings.ToLower(col), "ticket"):
			width = 18
		case strings.Contains(strings.ToLower(col), "description") || strings.Contains(strings.ToLower(col), "text") || strings.Contains(strings.ToLower(col), "solution"):
			width = 45
		case strings.Contains(strings.ToLower(col), "name"):
			width = 25
		case strings.Contains(strings.ToLower(col), "status"):
			width = 18
		case strings.Contains(strings.ToLower(col), "priority"):
			width = 12
		case strings.Contains(strings.ToLower(col), "date") || strings.Contains(strings.ToLower(col), "at"):
			width = 20
		case strings.Contains(strings.ToLower(col), "number"):
			width = 15
		default:
			width = 18
		}
		f.SetColWidth(tableName, colLetter, colLetter, width)
	}

	// Freeze first row
	f.SetPanes(tableName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Add Auto Filter
	if len(columns) > 0 {
		dataRange := fmt.Sprintf("A1:%s%d", string(rune('A'+len(columns)-1)), rowNum-1)
		f.AutoFilter(tableName, dataRange, []excelize.AutoFilterOptions{})
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
