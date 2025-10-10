package handlers

import (
	"fmt"
	"log"
	"reports-api/db"
	"reports-api/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get dashboard data
// @Description Get dashboard data including branches, departments, IP phones, programs, and tasks
// @Tags dashboard
// @Accept json
// @Produce json
// @Param month query string false "Month filter"
// @Param year query string false "Year filter"
// @Success 200 {object} models.DashboardResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/dashboard/data [get]
func GetDashboardDataHandler(c *fiber.Ctx) error {
	logger := log.Default()

	defer func() {
		if r := recover(); r != nil {
			logger.Printf("❌ PANIC in GetDashboardDataHandler: %v", r)
			c.Status(500).JSON(models.DashboardResponse{
				Success: false,
				Message: "Internal server error: " + fmt.Sprintf("%v", r),
			})
		}
	}()

	month := c.Query("month")
	year := c.Query("year")

	if db.DB == nil {
		logger.Printf("❌ ERROR: Database connection is nil")
		return c.Status(500).JSON(models.DashboardResponse{
			Success: false,
			Message: "Database connection error",
		})
	}

	err := db.DB.Ping()
	if err != nil {
		logger.Printf("❌ ERROR: Database ping failed: %v", err)
		return c.Status(500).JSON(models.DashboardResponse{
			Success: false,
			Message: "Database connection failed: " + err.Error(),
		})
	}

	// Fetch branches
	branches := []models.BranchDb{}
	branchRows, err := db.DB.Query("SELECT id, name, created_at, updated_at FROM branches WHERE deleted_at IS NULL")
	if err == nil {
		defer branchRows.Close()
		for branchRows.Next() {
			var b models.BranchDb
			err := branchRows.Scan(&b.ID, &b.Name, &b.CreatedAt, &b.UpdatedAt)
			if err == nil {
				branches = append(branches, b)
			}
		}
	}

	// Fetch departments
	departments := []models.DepartmentDb{}
	departmentRows, err := db.DB.Query(`
	SELECT d.id, d.name, d.branch_id, b.name as branch_name, d.created_at, d.updated_at
	FROM departments d
	LEFT JOIN branches b ON d.branch_id = b.id
	WHERE d.deleted_at IS NULL
  `)
	if err == nil {
		defer departmentRows.Close()
		for departmentRows.Next() {
			var d models.DepartmentDb
			err := departmentRows.Scan(
				&d.ID, &d.Name, &d.BranchID, &d.BranchName, &d.CreatedAt, &d.UpdatedAt,
			)
			if err == nil {
				scores := []models.Score{}
				scoresQuery := `SELECT department_id, year, month, score FROM scores WHERE department_id = ?`
				scoresRows, err := db.DB.Query(scoresQuery, d.ID)
				if err == nil {
					defer scoresRows.Close()
					for scoresRows.Next() {
						var s models.Score
						err := scoresRows.Scan(&s.DepartmentID, &s.Year, &s.Month, &s.Score)
						if err == nil {
							scores = append(scores, s)
						}
					}
				}

				total := 0
				for _, s := range scores {
					total += s.Score
				}
				d.Scores = total
				departments = append(departments, d)
			}
		}
	}

	// Fetch IP phones
	ipPhones := []models.IPPhoneDb{}
	ipPhoneRows, err := db.DB.Query(`
		SELECT 
			ip.id, ip.number, ip.name, ip.department_id, 
			d.name AS department_name, d.branch_id, 
			b.name AS branch_name, 
			ip.created_at, ip.updated_at
		FROM ip_phones ip
		LEFT JOIN departments d ON ip.department_id = d.id
		LEFT JOIN branches b ON d.branch_id = b.id
		WHERE ip.deleted_at IS NULL
	`)
	if err == nil {
		defer ipPhoneRows.Close()
		for ipPhoneRows.Next() {
			var ip models.IPPhoneDb
			var departmentName, branchName *string
			var branchID *int
			err := ipPhoneRows.Scan(
				&ip.ID, &ip.Number, &ip.Name, &ip.DepartmentID,
				&departmentName, &branchID, &branchName,
				&ip.CreatedAt, &ip.UpdatedAt,
			)
			if err == nil {
				ip.DepartmentName = departmentName
				ip.BranchID = branchID
				ip.BranchName = branchName
				ipPhones = append(ipPhones, ip)
			}
		}
	}

	// Fetch programs
	programs := []models.ProgramDb{}
	programRows, err := db.DB.Query("SELECT id, name, created_at, updated_at FROM systems_program WHERE deleted_at IS NULL")
	if err == nil {
		defer programRows.Close()
		for programRows.Next() {
			var p models.ProgramDb
			err := programRows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
			if err == nil {
				programs = append(programs, p)
			}
		}
	}

	// Fetch issue types
	issueTypes := []models.IssueTypeDb{}
	issueTypeRows, err := db.DB.Query("SELECT id, name, created_at FROM issue_types")
	if err == nil {
		defer issueTypeRows.Close()
		for issueTypeRows.Next() {
			var it models.IssueTypeDb
			err := issueTypeRows.Scan(&it.ID, &it.Name, &it.CreatedAt)
			if err == nil {
				issueTypes = append(issueTypes, it)
			}
		}
	}

	// Fetch tasks
	tasks := []models.TaskWithDetailsDb{}
	tasksQuery := `SELECT 
	       t.id, 
	       IFNULL(t.phone_id, 0), 
	       IFNULL(ip.number, 0), 
	       IFNULL(ip.name, ''), 
	       IFNULL(t.system_id, 0), 
	       IFNULL(sp.name, ''),
		   IFNULL(it.name, ''), 
	       IFNULL(t.department_id, 0), 
	       IFNULL(d.name, ''), 
	       IFNULL(d.branch_id, 0), 
	       IFNULL(b.name, ''),
		   IFNULL(t.reported_by, ''), 
	       IFNULL(t.text, ''), 
	       IFNULL(t.status, 0), 
	       IFNULL(t.created_at, ''), 
	       IFNULL(t.updated_at, '')
       FROM tasks t
       LEFT JOIN ip_phones ip ON t.phone_id = ip.id
       LEFT JOIN departments d ON t.department_id = d.id
       LEFT JOIN branches b ON d.branch_id = b.id
       LEFT JOIN systems_program sp ON t.system_id = sp.id
       LEFT JOIN issue_types it ON t.issue_type = it.id
       WHERE t.deleted_at IS NULL`
	var args []interface{}
	if month != "" && year != "" {
		tasksQuery += " AND MONTH(t.created_at) = ? AND YEAR(t.created_at) = ?"
		args = append(args, month, year)
	} else if month != "" {
		tasksQuery += " AND MONTH(t.created_at) = ?"
		args = append(args, month)
	} else if year != "" {
		tasksQuery += " AND YEAR(t.created_at) = ?"
		args = append(args, year)
	}
	tasksQuery += " ORDER BY t.created_at DESC"
	taskRows, err := db.DB.Query(tasksQuery, args...)
	if err == nil {
		defer taskRows.Close()
		for taskRows.Next() {
			var t models.TaskWithDetailsDb
			err := taskRows.Scan(&t.ID, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.SystemType, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.ReportedBy, &t.Text, &t.Status, &t.CreatedAt, &t.UpdatedAt)
			if err == nil {
				if t.CreatedAt != "" {
					parsed, err := time.Parse("2006-01-02 15:04:05", t.CreatedAt)
					if err == nil {
						t.Month = parsed.Month().String()
						t.Year = fmt.Sprintf("%d", parsed.Year())
					}
				}
				tasks = append(tasks, t)
			}
		}
	}

	chartData := calculateChartData(tasks)

	response := models.DashboardResponse{
		Success:     true,
		Message:     "Dashboard data retrieved successfully",
		ChartData:   chartData,
		Branches:    branches,
		Departments: departments,
		IPPhones:    ipPhones,
		Programs:    programs,
		Tasks:       tasks,
		IssueTypes:  issueTypes,
	}
	return c.JSON(response)
}

// calculateChartData calculates data for charts
func calculateChartData(tasks []models.TaskWithDetailsDb) models.ChartData {
	return models.ChartData{
		YearStats: []models.YearStat{},
	}
}
