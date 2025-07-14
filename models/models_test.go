package models

import (
	"database/sql"
	"testing"
)

func TestNullStringToPtr(t *testing.T) {
	ns := sql.NullString{String: "test", Valid: true}
	ptr := nullStringToPtr(ns)
	if ptr == nil || *ptr != "test" {
		t.Errorf("Expected pointer to 'test', got %v", ptr)
	}

	ns = sql.NullString{String: "", Valid: false}
	ptr = nullStringToPtr(ns)
	if ptr != nil {
		t.Errorf("Expected nil, got %v", ptr)
	}
}

func TestNullInt64ToPtr(t *testing.T) {
	n := sql.NullInt64{Int64: 42, Valid: true}
	ptr := nullInt64ToPtr(n)
	if ptr == nil || *ptr != 42 {
		t.Errorf("Expected pointer to 42, got %v", ptr)
	}

	n = sql.NullInt64{Int64: 0, Valid: false}
	ptr = nullInt64ToPtr(n)
	if ptr != nil {
		t.Errorf("Expected nil, got %v", ptr)
	}
}

func TestStructInitialization(t *testing.T) {
	user := User{ID: 1, Username: sql.NullString{String: "admin", Valid: true}, Role: sql.NullString{String: "admin", Valid: true}}
	if user.ID != 1 || !user.Username.Valid || user.Username.String != "admin" {
		t.Errorf("User struct not initialized correctly: %+v", user)
	}
}
