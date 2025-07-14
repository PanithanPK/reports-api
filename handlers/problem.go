package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

// Custom logger for problem handlers
var problemLogger = log.New(log.Writer(), "[PROBLEM] ", log.Ldate|log.Ltime)

// ReportProblemHandler handles POST requests to create new problem reports
func ReportProblemHandler(w http.ResponseWriter, r *http.Request) {
	problemLogger.Printf("📝 Creating new problem report from %s", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	var req models.ReportProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problemLogger.Printf("❌ Invalid request body from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	problemLogger.Printf("📋 Problem report data - IP: %s, Program: %s, Other: %s",
		req.IpPhone, req.Program, req.Other)

	// Validate required fields
	if (req.IpPhone == "" && req.Other == "") || req.Program == "" || req.Problem == "" {
		problemLogger.Printf("❌ Missing required fields from %s", r.RemoteAddr)
		response := models.ReportProblemResponse{
			Success: false,
			Message: "Missing required fields: ipPhone, program, problem",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Insert into database
	query := `INSERT INTO report_problem (ip_phone, program, other, problem, status) VALUES (NULLIF(?, ''), ?, ?, ?, 'Pending')`
	result, err := db.DB.Exec(query, req.IpPhone, req.Program, req.Other, req.Problem, req.status)
	if err != nil {
		problemLogger.Printf("❌ Error inserting problem: %v", err)
		response := models.ReportProblemResponse{
			Success: false,
			Message: "Failed to create problem report",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get the inserted ID
	id, _ := result.LastInsertId()
	problemLogger.Printf("✅ Problem report created successfully with ID: %d", id)

	response := models.ReportProblemResponse{
		Success: true,
		Message: "Problem report created successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// GetProblemsHandler handles GET requests to retrieve all problems
func GetProblemsHandler(w http.ResponseWriter, r *http.Request) {
	problemLogger.Printf("📋 Retrieving all problems from %s", r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT rp.id, rp.ip_phone, rp.program, rp.other, rp.problem, rp.solution, 
		       rp.solution_date, rp.solution_user, rp.status, rp.created_at,
		       bo.branchoffice
		FROM report_problem rp
		LEFT JOIN branch_office bo ON rp.ip_phone = bo.ip_phone
		ORDER BY rp.created_at DESC
	`

	rows, err := db.DB.Query(query)
	if err != nil {
		problemLogger.Printf("❌ Error querying problems: %v", err)
		response := models.GetProblemsResponse{
			Success: false,
			Message: "Failed to retrieve problems",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()

	var problems []models.Problem
	var scanErrors int

	for rows.Next() {
		var p models.Problem
		var solutionDateStr, createdAtStr sql.NullString

		err := rows.Scan(
			&p.ID, &p.IpPhone, &p.Program, &p.Other, &p.Problem, &p.Solution,
			&solutionDateStr, &p.SolutionUser, &p.Status, &createdAtStr, &p.Branchoffice,
		)
		if err != nil {
			problemLogger.Printf("❌ Error scanning problem: %v", err)
			scanErrors++
			continue
		}

		// Parse solution_date if valid
		if solutionDateStr.Valid {
			t, err := time.Parse("2006-01-02 15:04:05", solutionDateStr.String)
			if err != nil {
				problemLogger.Printf("⚠️ Error parsing solution_date for problem %d: %v", p.ID, err)
			} else {
				p.SolutionDate = sql.NullTime{Time: t, Valid: true}
			}
		}

		// Parse created_at if valid
		if createdAtStr.Valid {
			t, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
			if err != nil {
				problemLogger.Printf("⚠️ Error parsing created_at for problem %d: %v", p.ID, err)
			} else {
				p.CreatedAt = sql.NullTime{Time: t, Valid: true}
			}
		}

		problems = append(problems, p)
	}

	if scanErrors > 0 {
		problemLogger.Printf("⚠️ %d problems had scanning errors", scanErrors)
	}

	problemLogger.Printf("✅ Retrieved %d problems successfully", len(problems))

	response := models.GetProblemsResponse{
		Success:  true,
		Message:  "Problems retrieved successfully",
		Problems: problems,
	}
	json.NewEncoder(w).Encode(response)
}

// GetProblemByIDHandler handles GET requests to retrieve a specific problem by ID
func GetProblemByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["id"]

	problemLogger.Printf("🔍 Retrieving problem ID: %s from %s", problemID, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	query := `
		SELECT rp.id, rp.ip_phone, rp.program, rp.other, rp.problem, rp.solution, 
		       rp.solution_date, rp.solution_user, rp.status, rp.created_at,
		       COALESCE(bo.branchoffice, '') as branchoffice
		FROM report_problem rp
		LEFT JOIN branch_office bo ON rp.ip_phone = bo.ip_phone
		WHERE rp.id = ?
	`

	var p models.Problem
	var solutionDateStr, createdAtStr sql.NullString

	err := db.DB.QueryRow(query, problemID).Scan(
		&p.ID, &p.IpPhone, &p.Program, &p.Other, &p.Problem, &p.Solution,
		&solutionDateStr, &p.SolutionUser, &p.Status, &createdAtStr, &p.Branchoffice,
	)

	if err == sql.ErrNoRows {
		problemLogger.Printf("❌ Problem ID %s not found", problemID)
		response := models.GetProblemResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err != nil {
		problemLogger.Printf("❌ Error querying problem ID %s: %v", problemID, err)
		response := models.GetProblemResponse{
			Success: false,
			Message: "Failed to retrieve problem",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse solution_date if valid
	if solutionDateStr.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", solutionDateStr.String)
		if err != nil {
			problemLogger.Printf("⚠️ Error parsing solution_date for problem %d: %v", p.ID, err)
		} else {
			p.SolutionDate = sql.NullTime{Time: t, Valid: true}
		}
	}

	// Parse created_at if valid
	if createdAtStr.Valid {
		t, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String)
		if err != nil {
			problemLogger.Printf("⚠️ Error parsing created_at for problem %d: %v", p.ID, err)
		} else {
			p.CreatedAt = sql.NullTime{Time: t, Valid: true}
		}
	}

	problemLogger.Printf("✅ Problem ID %s retrieved successfully", problemID)

	response := models.GetProblemResponse{
		Success: true,
		Message: "Problem retrieved successfully",
		Problem: p,
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateProblemHandler handles PUT requests to update existing problems
func UpdateProblemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["id"]

	problemLogger.Printf("✏️ Updating problem ID: %s from %s", problemID, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	var req models.UpdateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problemLogger.Printf("❌ Invalid request body for problem ID %s: %v", problemID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	problemLogger.Printf("📋 Update data for problem %s - IP: %s, Program: %s, Other: %s",
		problemID, req.IpPhone, req.Program, req.Other)

	// Update database
	query := `UPDATE report_problem SET ip_phone = ?, other = ?, program = ?, problem = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	result, err := db.DB.Exec(query, req.IpPhone, req.Other, req.Program, req.Problem, problemID)
	if err != nil {
		problemLogger.Printf("❌ Error updating problem ID %s: %v", problemID, err)
		response := models.UpdateProblemResponse{
			Success: false,
			Message: "Failed to update problem",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		problemLogger.Printf("❌ Problem ID %s not found for update", problemID)
		response := models.UpdateProblemResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	problemLogger.Printf("✅ Problem ID %s updated successfully (%d rows affected)", problemID, rowsAffected)

	response := models.UpdateProblemResponse{
		Success: true,
		Message: "Problem updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteProblemHandler handles DELETE requests to delete problems
func DeleteProblemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	problemID := vars["id"]

	problemLogger.Printf("🗑️ Deleting problem ID: %s from %s", problemID, r.RemoteAddr)

	w.Header().Set("Content-Type", "application/json")

	// Delete from database
	query := `DELETE FROM report_problem WHERE id = ?`
	result, err := db.DB.Exec(query, problemID)
	if err != nil {
		problemLogger.Printf("❌ Error deleting problem ID %s: %v", problemID, err)
		response := models.DeleteProblemResponse{
			Success: false,
			Message: "Failed to delete problem",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		problemLogger.Printf("❌ Problem ID %s not found for deletion", problemID)
		response := models.DeleteProblemResponse{
			Success: false,
			Message: "Problem not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	problemLogger.Printf("✅ Problem ID %s deleted successfully (%d rows affected)", problemID, rowsAffected)

	response := models.DeleteProblemResponse{
		Success: true,
		Message: "Problem deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteAllProblems ลบข้อมูลปัญหาทั้งหมด
func DeleteAllProblemsHandler(w http.ResponseWriter, r *http.Request) {
	problemLogger.Printf("🗑️ Deleting all problems from %s", r.RemoteAddr)

	database := db.DB

	// ลบข้อมูลทั้งหมดจากตาราง report_problem และรีเซ็ต auto-increment
	// ใช้ DELETE FROM แทน TRUNCATE TABLE เพื่อหลีกเลี่ยง foreign key constraints
	result, err := database.Exec("DELETE FROM report_problem")
	if err != nil {
		problemLogger.Printf("❌ Error deleting all problems: %v", err)
		response := models.DeleteAllProblemsResponse{
			Success: false,
			Message: "Failed to delete all problems: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	problemLogger.Printf("✅ Deleted %d problems successfully", rowsAffected)

	// รีเซ็ต auto-increment counter
	_, err = database.Exec("ALTER TABLE report_problem AUTO_INCREMENT = 1")
	if err != nil {
		problemLogger.Printf("⚠️ Error resetting auto increment: %v", err)
		// ไม่ return error เพราะข้อมูลถูกลบแล้ว
	} else {
		problemLogger.Printf("✅ Auto increment reset successfully")
	}

	response := models.DeleteAllProblemsResponse{
		Success: true,
		Message: "All problems deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// GetDashboardDataHandler ดึงข้อมูลสำหรับ Dashboard
func GetDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	problemLogger.Printf("📊 Loading dashboard data from %s", r.RemoteAddr)

	defer func() {
		if r := recover(); r != nil {
			problemLogger.Printf("❌ PANIC in GetDashboardDataHandler: %v", r)
			response := models.DashboardResponse{
				Success: false,
				Message: "Internal server error: " + fmt.Sprintf("%v", r),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
	}()
	w.Header().Set("Content-Type", "application/json")

	// อ่าน query parameter month และ year
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	if month != "" || year != "" {
		problemLogger.Printf("📅 Dashboard filter - Month: %s, Year: %s", month, year)
	}

	if db.DB == nil {
		problemLogger.Printf("❌ ERROR: Database connection is nil")
		response := models.DashboardResponse{
			Success: false,
			Message: "Database connection error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := db.DB.Ping()
	if err != nil {
		problemLogger.Printf("❌ ERROR: Database ping failed: %v", err)
		response := models.DashboardResponse{
			Success: false,
			Message: "Database connection failed: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	problemLogger.Printf("✅ Database connection verified")

	// ปรับ query ให้รองรับ filter เดือนและปี
	problemsQuery := `
		SELECT p.id, p.ip_phone, p.program, p.other, p.problem, p.solution, 
		       p.solution_date, p.solution_user, p.status, p.created_at
		FROM report_problem p
		WHERE 1=1
	`
	var args []interface{}
	if month != "" && year != "" {
		problemsQuery += " AND MONTH(p.created_at) = ? AND YEAR(p.created_at) = ?"
		args = append(args, month, year)
	} else if month != "" {
		problemsQuery += " AND MONTH(p.created_at) = ?"
		args = append(args, month)
	} else if year != "" {
		problemsQuery += " AND YEAR(p.created_at) = ?"
		args = append(args, year)
	}
	problemsQuery += " ORDER BY p.created_at DESC"

	problemRows, err := db.DB.Query(problemsQuery, args...)
	if err != nil {
		log.Printf("ERROR querying problems: %v", err)
		response := models.DashboardResponse{
			Success: false,
			Message: "Failed to retrieve problems data: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer problemRows.Close()

	var problems []models.Problem
	for problemRows.Next() {
		var p models.Problem
		var createdAtStr, solutionDateStr sql.NullString
		err := problemRows.Scan(
			&p.ID, &p.IpPhone, &p.Program, &p.Other, &p.Problem,
			&p.Solution, &solutionDateStr, &p.SolutionUser, &p.Status, &createdAtStr,
		)
		if err != nil {
			log.Printf("Error scanning problem: %v", err)
			continue
		}

		// Parse created_at และคำนวณ month, year
		if createdAtStr.Valid && createdAtStr.String != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr.String); err == nil {
				p.CreatedAt = sql.NullTime{Time: t, Valid: true}

				// คำนวณ month และ year จาก created_at
				monthName := t.Month().String()
				yearStr := fmt.Sprintf("%d", t.Year())

				p.Month = sql.NullString{String: monthName, Valid: true}
				p.Year = sql.NullString{String: yearStr, Valid: true}
			} else {
				p.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
				// ถ้า parse ไม่ได้ ใช้เวลาปัจจุบัน
				now := time.Now()
				p.Month = sql.NullString{String: now.Month().String(), Valid: true}
				p.Year = sql.NullString{String: fmt.Sprintf("%d", now.Year()), Valid: true}
			}
		} else {
			// ถ้าไม่มี created_at ใช้เวลาปัจจุบัน
			now := time.Now()
			p.CreatedAt = sql.NullTime{Time: now, Valid: true}
			p.Month = sql.NullString{String: now.Month().String(), Valid: true}
			p.Year = sql.NullString{String: fmt.Sprintf("%d", now.Year()), Valid: true}
		}

		// Parse solution_date
		if solutionDateStr.Valid && solutionDateStr.String != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", solutionDateStr.String); err == nil {
				p.SolutionDate = sql.NullTime{Time: t, Valid: true}
			}
		}

		problems = append(problems, p)
	}

	branchQuery := `SELECT ip_phone, branchoffice FROM branch_office ORDER BY branchoffice`
	branchRows, err := db.DB.Query(branchQuery)
	if err != nil {
		log.Printf("ERROR querying branch offices: %v", err)
		response := models.DashboardResponse{
			Success: false,
			Message: "Failed to retrieve branch offices data: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer branchRows.Close()

	var branchOffices []models.BranchOffice
	for branchRows.Next() {
		var bo models.BranchOffice
		err := branchRows.Scan(&bo.IpPhone, &bo.Branchoffice)
		if err != nil {
			log.Printf("Error scanning branch office: %v", err)
			continue
		}
		branchOffices = append(branchOffices, bo)
	}

	programQuery := `SELECT id, name FROM program ORDER BY name`
	programRows, err := db.DB.Query(programQuery)
	if err != nil {
		log.Printf("ERROR querying programs: %v", err)
		response := models.DashboardResponse{
			Success: false,
			Message: "Failed to retrieve programs data: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	defer programRows.Close()

	var programs []models.Program
	for programRows.Next() {
		var program models.Program
		err := programRows.Scan(&program.ID, &program.Name)
		if err != nil {
			log.Printf("Error scanning program: %v", err)
			continue
		}
		programs = append(programs, program)
	}

	chartData := calculateChartData(problems, branchOffices, programs)

	err = json.NewEncoder(w).Encode(models.DashboardResponse{
		Success:       true,
		Message:       "Dashboard data retrieved successfully",
		Problems:      problems,
		BranchOffices: branchOffices,
		Programs:      programs,
		ChartData:     chartData,
	})
	if err != nil {
		log.Printf("ERROR encoding response: %v", err)
		return
	}
}

// calculateChartData คำนวณข้อมูลสำหรับกราฟ
func calculateChartData(problems []models.Problem, branchOffices []models.BranchOffice, programs []models.Program) models.ChartData {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in calculateChartData: %v", r)
		}
	}()

	log.Printf("Starting chart calculation with %d problems, %d branches, %d programs",
		len(problems), len(branchOffices), len(programs))

	branchStats := make(map[string]int)
	programStats := make(map[string]int)
	branchProgramStats := make(map[string]map[string]int)

	for _, problem := range problems {
		branchName := "ไม่ระบุ"
		found := false
		if problem.IpPhone.Valid && problem.IpPhone.String != "" {
			for _, branch := range branchOffices {
				if branch.IpPhone.Valid && branch.IpPhone.String == problem.IpPhone.String {
					if branch.Branchoffice.Valid {
						branchName = branch.Branchoffice.String
					} else {
						branchName = branch.IpPhone.String
					}
					found = true
					break
				}
			}
			if !found {
				branchName = "อื่นๆ"
			}
		} else if problem.Other.Valid && problem.Other.String != "" {
			branchName = "อื่นๆ"
		}

		programName := "ไม่ระบุ"
		if problem.Program.Valid && problem.Program.String != "" {
			programName = problem.Program.String
		}

		branchStats[branchName]++
		programStats[programName]++

		if branchProgramStats[branchName] == nil {
			branchProgramStats[branchName] = make(map[string]int)
		}
		branchProgramStats[branchName][programName]++
	}

	var barChartData []models.BarChartData
	for branchName, totalCount := range branchStats {
		programData := make(map[string]int)
		if branchProgramStats[branchName] != nil {
			programData = branchProgramStats[branchName]
		}
		barChartData = append(barChartData, models.BarChartData{
			BranchName:  branchName,
			TotalCount:  totalCount,
			ProgramData: programData,
		})
	}

	sort.Slice(barChartData, func(i, j int) bool {
		return barChartData[i].TotalCount > barChartData[j].TotalCount
	})

	var pieChartData []models.PieChartData
	for programName, count := range programStats {
		percentage := 0.0
		if len(problems) > 0 {
			percentage = float64(count) / float64(len(problems)) * 100
		}
		pieChartData = append(pieChartData, models.PieChartData{
			ProgramName: programName,
			Count:       count,
			Percentage:  percentage,
		})
	}

	sort.Slice(pieChartData, func(i, j int) bool {
		return pieChartData[i].Count > pieChartData[j].Count
	})

	return models.ChartData{
		BarChartData: barChartData,
		PieChartData: pieChartData,
	}
}
