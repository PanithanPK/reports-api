package handlers

import (
	"encoding/json"
	"net/http"
	"reports-api/backup"
	"strconv"
)

func BackupHandler(w http.ResponseWriter, r *http.Request) {
	bs := backup.NewBackupService()
	
	if err := bs.CreateBackup(); err != nil {
		http.Error(w, "Backup failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"message": "Backup created successfully"})
}

func CleanBackupsHandler(w http.ResponseWriter, r *http.Request) {
	daysStr := r.URL.Query().Get("days")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7 // default 7 days
	}
	
	bs := backup.NewBackupService()
	if err := bs.CleanOldBackups(days); err != nil {
		http.Error(w, "Clean failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(map[string]string{"message": "Old backups cleaned"})
}