package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reports-api/models"
)

// SendTelegramNotificationHandler รับ POST แล้วส่งข้อความไป Telegram
func SendTelegramNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.TelegramResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req models.TelegramRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Reportmessage == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.TelegramResponse{Success: false, Message: "Invalid request body"})
		return
	}

	// สร้างข้อความรวมข้อมูล
	msg := ""
	if req.BranchName != "" {
		msg += "สาขา: " + req.BranchName + "\n"
	}
	if req.DepartmentName != "" {
		msg += "แผนก: " + req.DepartmentName + "\n"
	}
	if req.Number != "" {
		msg += "เบอร์: " + req.Number + "\n"
	}
	if req.IPPhoneName != "" {
		msg += "ผู้รับผิดชอบ: " + req.IPPhoneName + "\n"
	}
	if req.Reportmessage != "" {
		msg += "รายงานปัญหา: " + req.Reportmessage + "\n"
	}
	if req.URL != "" {
		msg += "[ดูรายละเอียดเพิ่มเติม](" + req.URL + ")\n"
	}

	botToken := os.Getenv("botToken")
	chatID := os.Getenv("chatID")

	// แสดงสภาพแวดล้อมที่กำลังใช้งาน
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "default"
	}

	fmt.Printf("[Telegram][%s], chatID: %s\n", env, chatID)

	telegramAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := map[string]interface{}{
		"chat_id":                  chatID,
		"text":                     msg,
		"parse_mode":               "Markdown",
		"disable_web_page_preview": false,
	}
	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(telegramAPI, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.TelegramResponse{Success: false, Message: "Failed to send telegram message: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(models.TelegramResponse{Success: false, Message: "Telegram API error: " + string(body)})
		return
	}

	json.NewEncoder(w).Encode(models.TelegramResponse{Success: true, Message: "Notification sent successfully"})
}
