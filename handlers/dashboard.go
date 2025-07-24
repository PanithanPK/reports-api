package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reports-api/db"
	"reports-api/models"
	"time"
)

// GetDashboardDataHandler ดึงข้อมูลสำหรับ Dashboard
func GetDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Default()
	logger.Printf("📊 Loading dashboard data from %s", r.RemoteAddr)

	defer func() {
		if r := recover(); r != nil {
			logger.Printf("❌ PANIC in GetDashboardDataHandler: %v", r)
			response := models.DashboardResponse{
				Success: false,
				Message: "Internal server error: " + fmt.Sprintf("%v", r),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
	}()

	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	if db.DB == nil {
		logger.Printf("❌ ERROR: Database connection is nil")
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
		logger.Printf("❌ ERROR: Database ping failed: %v", err)
		response := models.DashboardResponse{
			Success: false,
			Message: "Database connection failed: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// ดึงข้อมูล branches
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

	// ดึงข้อมูล departments
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
				// ดึงข้อมูล scores สำหรับแต่ละ department
				scores := []models.Score{} // ประกาศเป็น slice ที่ว่าง
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

				// เพิ่ม scores เข้าไปใน department
				total := 0
				for _, s := range scores {
					total += s.Score
				}
				d.Scores = total
				departments = append(departments, d)
			} else {
				log.Println("Scan error (departments):", err)
			}
		}
	}

	// ดึงข้อมูล ip_phones
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

	// ดึงข้อมูล programs
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

	// ดึงข้อมูล tasks พร้อม join ip_phones, departments, branches, systems_program
	tasks := []models.TaskWithDetailsDb{}
	tasksQuery := `SELECT t.id, t.phone_id, ip.number, ip.name, t.system_id, sp.name, ip.department_id, d.name, d.branch_id, b.name, t.text, t.status, t.created_at, t.updated_at
	FROM tasks t
	LEFT JOIN ip_phones ip ON t.phone_id = ip.id
	LEFT JOIN departments d ON ip.department_id = d.id
	LEFT JOIN branches b ON d.branch_id = b.id
	LEFT JOIN systems_program sp ON t.system_id = sp.id
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
			err := taskRows.Scan(&t.ID, &t.PhoneID, &t.Number, &t.PhoneName, &t.SystemID, &t.SystemName, &t.DepartmentID, &t.DepartmentName, &t.BranchID, &t.BranchName, &t.Text, &t.Status, &t.CreatedAt, &t.UpdatedAt)
			if err == nil {
				// แปลง t.CreatedAt (string) เป็น time.Time เพื่อดึง month/year
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
	}
	json.NewEncoder(w).Encode(response)
}

// calculateChartData คำนวณข้อมูลสำหรับกราฟ
func calculateChartData(tasks []models.TaskWithDetailsDb) models.ChartData {
	yearMonthBranchDeptIPProg := make(map[string]map[string]map[string]map[string]map[string]map[string]int)
	type idMap struct {
		branchId     *int
		departmentId *int
		ipphoneId    *int
		programId    *int
	}
	idMaps := make(map[string]map[string]map[string]map[string]map[string]idMap)
	// year -> month -> branch -> department -> ipphone -> idMap

	for _, t := range tasks {
		year := t.Year
		month := t.Month
		branch := "ไม่ระบุ"
		if t.BranchName != nil && *t.BranchName != "" {
			branch = *t.BranchName
		}
		department := "ไม่ระบุ"
		if t.DepartmentName != nil && *t.DepartmentName != "" {
			department = *t.DepartmentName
		}
		ipphone := "ไม่ระบุ"
		if t.PhoneName != nil && *t.PhoneName != "" {
			ipphone = *t.PhoneName
		}
		program := "ไม่ระบุ"
		if t.SystemName != nil && *t.SystemName != "" {
			program = *t.SystemName
		}
		if yearMonthBranchDeptIPProg[year] == nil {
			yearMonthBranchDeptIPProg[year] = make(map[string]map[string]map[string]map[string]map[string]int)
		}
		if yearMonthBranchDeptIPProg[year][month] == nil {
			yearMonthBranchDeptIPProg[year][month] = make(map[string]map[string]map[string]map[string]int)
		}
		if yearMonthBranchDeptIPProg[year][month][branch] == nil {
			yearMonthBranchDeptIPProg[year][month][branch] = make(map[string]map[string]map[string]int)
		}
		if yearMonthBranchDeptIPProg[year][month][branch][department] == nil {
			yearMonthBranchDeptIPProg[year][month][branch][department] = make(map[string]map[string]int)
		}
		if yearMonthBranchDeptIPProg[year][month][branch][department][ipphone] == nil {
			yearMonthBranchDeptIPProg[year][month][branch][department][ipphone] = make(map[string]int)
		}
		yearMonthBranchDeptIPProg[year][month][branch][department][ipphone][program]++

		// map id
		if idMaps[year] == nil {
			idMaps[year] = make(map[string]map[string]map[string]map[string]idMap)
		}
		if idMaps[year][month] == nil {
			idMaps[year][month] = make(map[string]map[string]map[string]idMap)
		}
		if idMaps[year][month][branch] == nil {
			idMaps[year][month][branch] = make(map[string]map[string]idMap)
		}
		if idMaps[year][month][branch][department] == nil {
			idMaps[year][month][branch][department] = make(map[string]idMap)
		}
		idMaps[year][month][branch][department][ipphone] = idMap{
			branchId:     t.BranchID,
			departmentId: t.DepartmentID,
			ipphoneId:    &t.PhoneID,
			programId:    t.SystemID,
		}
	}

	var yearStats []models.YearStat
	for year, months := range yearMonthBranchDeptIPProg {
		var monthStats []models.MonthStat
		for month, branches := range months {
			var branchStats []models.BranchStat
			for branch, departments := range branches {
				var departmentStats []models.DepartmentStat
				totalBranchProblems := 0
				for department, ipphones := range departments {
					var ipphoneStats []models.IPPhoneStat
					totalDeptProblems := 0
					for ipphone, programs := range ipphones {
						var programStats []models.ProgramStat
						totalIPProblems := 0
						for program, count := range programs {
							id := idMaps[year][month][branch][department][ipphone]
							programStats = append(programStats, models.ProgramStat{
								ProgramId:     id.programId,
								ProgramName:   program,
								TotalProblems: count,
							})
							totalIPProblems += count
						}
						id := idMaps[year][month][branch][department][ipphone]
						ipphoneStats = append(ipphoneStats, models.IPPhoneStat{
							IPPhoneId:     id.ipphoneId,
							IPPhoneName:   ipphone,
							TotalProblems: totalIPProblems,
							Programs:      programStats,
						})
						totalDeptProblems += totalIPProblems
					}
					var deptId *int
					for _, v := range idMaps[year][month][branch][department] {
						deptId = v.departmentId
						break
					}
					departmentStats = append(departmentStats, models.DepartmentStat{
						DepartmentId:   deptId,
						DepartmentName: department,
						TotalProblems:  totalDeptProblems,
						IPPhones:       ipphoneStats,
					})
					totalBranchProblems += totalDeptProblems
				}
				var branchId *int
				for _, deptMap := range idMaps[year][month][branch] {
					for _, v := range deptMap {
						branchId = v.branchId
						break
					}
					if branchId != nil {
						break
					}
				}
				branchStats = append(branchStats, models.BranchStat{
					BranchId:      branchId,
					BranchName:    branch,
					TotalProblems: totalBranchProblems,
					Departments:   departmentStats,
				})
			}
			monthStats = append(monthStats, models.MonthStat{
				Month:    month,
				Branches: branchStats,
			})
		}
		yearStats = append(yearStats, models.YearStat{
			Year:   year,
			Months: monthStats,
		})
	}

	return models.ChartData{
		YearStats: yearStats,
	}
}
