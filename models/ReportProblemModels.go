package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// ReportProblemRequest ใช้สำหรับรับข้อมูลจาก client
type ReportProblemRequest struct {
	IpPhone string `json:"ipPhone"`
	Other   string `json:"other"`
	Program string `json:"program"`
	Problem string `json:"problem"`
	Status  string `json:"status"`
}

// ReportProblemResponse ตอบกลับการรายงานปัญหา
type ReportProblemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SolveProblemRequest รับข้อมูลแก้ไขปัญหา
type SolveProblemRequest struct {
	ProblemID    int          `json:"problemId"`
	Solution     string       `json:"solution"`
	SolutionDate sql.NullTime `json:"solutionDate,omitempty"`
	SolutionUser string       `json:"solutionUser"`
}

// SolveProblemResponse ตอบกลับการแก้ไขปัญหา
type SolveProblemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ResetSolutionResponse ตอบกลับการรีเซ็ตข้อมูลการแก้ไขปัญหา
type ResetSolutionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Problem แทนข้อมูลปัญหา รายงาน รวม solution
type Problem struct {
	ID           int            `json:"id"`
	IpPhone      sql.NullString `json:"ipPhone"`
	Program      sql.NullString `json:"program"`
	Other        sql.NullString `json:"other"`
	Problem      sql.NullString `json:"problem"`
	Solution     sql.NullString `json:"solution"`
	SolutionDate sql.NullTime   `json:"-"`
	SolutionUser sql.NullString `json:"solutionUser"`
	Status       sql.NullString `json:"status"`
	CreatedAt    sql.NullTime   `json:"-"`
	Branchoffice sql.NullString `json:"branchoffice"`
	Month        sql.NullString `json:"month"`
	Year         sql.NullString `json:"year"`
}

// MarshalJSON แปลง sql.NullString, sql.NullTime เป็น format ที่ frontend ต้องการ
func (p *Problem) MarshalJSON() ([]byte, error) {
	// สร้าง struct สำหรับ sql.NullString ที่มี String field
	type NullStringWrapper struct {
		String string `json:"String"`
		Valid  bool   `json:"Valid"`
	}

	formatNullString := func(ns sql.NullString) interface{} {
		if ns.Valid && ns.String != "" {
			return NullStringWrapper{String: ns.String, Valid: true}
		}
		return NullStringWrapper{String: "", Valid: false}
	}

	formatNullTime := func(nt sql.NullTime) *string {
		if nt.Valid {
			str := nt.Time.Format(time.RFC3339)
			return &str
		}
		return nil
	}

	type Alias struct {
		ID           int         `json:"id"`
		IpPhone      interface{} `json:"ipPhone"`
		Program      interface{} `json:"program"`
		Other        interface{} `json:"other"`
		Problem      interface{} `json:"problem"`
		Solution     interface{} `json:"solution"`
		SolutionDate *string     `json:"solutionDate"`
		SolutionUser interface{} `json:"solutionUser"`
		Status       interface{} `json:"status"`
		CreatedAt    *string     `json:"createdAt"`
		Branchoffice interface{} `json:"branchoffice"`
		Month        interface{} `json:"month"`
		Year         interface{} `json:"year"`
	}

	return json.Marshal(&Alias{
		ID:           p.ID,
		IpPhone:      formatNullString(p.IpPhone),
		Program:      formatNullString(p.Program),
		Other:        formatNullString(p.Other),
		Problem:      formatNullString(p.Problem),
		Solution:     formatNullString(p.Solution),
		SolutionDate: formatNullTime(p.SolutionDate),
		SolutionUser: formatNullString(p.SolutionUser),
		Status:       formatNullString(p.Status),
		CreatedAt:    formatNullTime(p.CreatedAt),
		Branchoffice: formatNullString(p.Branchoffice),
		Month:        formatNullString(p.Month),
		Year:         formatNullString(p.Year),
	})
}

// GetProblemsResponse โครงสร้างตอบกลับรายงานปัญหาแบบหลายรายการ
type GetProblemsResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message"`
	Problems []Problem `json:"problems"`
}

// GetProblemResponse โครงสร้างตอบกลับรายงานปัญหารายการเดียว
type GetProblemResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Problem Problem `json:"problem"`
}

// UserRequest โครงสร้างรับข้อมูลผู้ใช้
type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

// UserResponse ตอบกลับผู้ใช้
type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// User โครงสร้างข้อมูลผู้ใช้
type User struct {
	ID       int            `json:"id"`
	Username sql.NullString `json:"username"`
	Role     sql.NullString `json:"role"`
}

// GetUsersResponse โครงสร้างตอบกลับรายชื่อผู้ใช้
type GetUsersResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Users   []User `json:"users"`
}

// UpdateProblemRequest สำหรับอัพเดตข้อมูลปัญหา
type UpdateProblemRequest struct {
	ID      int    `json:"id"`
	IpPhone string `json:"ipPhone"`
	Other   string `json:"other"`
	Program string `json:"program"`
	Problem string `json:"problem"`
}

// UpdateProblemResponse ตอบกลับการอัพเดตปัญหา
type UpdateProblemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DeleteProblemResponse ตอบกลับการลบปัญหา
type DeleteProblemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateSolutionRequest รับข้อมูลการแก้ไขข้อมูลการแก้ไขปัญหา
type UpdateSolutionRequest struct {
	Solution     string `json:"solution"`
	SolutionUser string `json:"solutionUser"`
}

// UpdateSolutionResponse ตอบกลับการแก้ไขข้อมูลการแก้ไขปัญหา
type UpdateSolutionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DeleteAllProblemsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Dashboard Models
type DashboardResponse struct {
	Success       bool           `json:"success"`
	Message       string         `json:"message"`
	Problems      []Problem      `json:"problems"`
	BranchOffices []BranchOffice `json:"branchOffices"`
	Programs      []Program      `json:"programs"`
	ChartData     ChartData      `json:"chartData"`
}

type ChartData struct {
	BarChartData []BarChartData `json:"barChartData"`
	PieChartData []PieChartData `json:"pieChartData"`
}

type BarChartData struct {
	BranchName  string         `json:"branchName"`
	TotalCount  int            `json:"totalCount"`
	ProgramData map[string]int `json:"programData"`
}

type PieChartData struct {
	ProgramName string  `json:"programName"`
	Count       int     `json:"count"`
	Percentage  float64 `json:"percentage"`
}

type UpdateProblemANDSolveProblemRequest struct {
	Problem      string `json:"problem"`
	Solution     string `json:"solution"`
	SolutionUser string `json:"solutionUser"`
}

type UpdateProblemANDSolveProblemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DeleteProblemandSolveProblemResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
