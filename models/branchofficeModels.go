package models

import (
	"database/sql"
	"encoding/json"
)

type BranchOfficeRequest struct {
	IpPhone      string `json:"ipPhone"`
	Branchoffice string `json:"branchoffice"`
}

type BranchOfficeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GetBranchOfficesResponse struct {
	Success       bool           `json:"success"`
	Message       string         `json:"message"`
	BranchOffices []BranchOffice `json:"branchOffices"`
	Programs      []Program      `json:"programs"`
}

// BranchOffice โครงสร้างข้อมูลสาขา
type BranchOffice struct {
	IpPhone      sql.NullString `json:"-"`
	Branchoffice sql.NullString `json:"-"`
}

// MarshalJSON แปลง sql.NullString เป็น string หรือ null JSON
func (bo *BranchOffice) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		IpPhone      *string `json:"ip_phone"`
		Branchoffice *string `json:"branchoffice"`
	}{
		IpPhone:      nullStringToPtr(bo.IpPhone),
		Branchoffice: nullStringToPtr(bo.Branchoffice),
	})
}

type DeleteAllBranchOfficesResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
