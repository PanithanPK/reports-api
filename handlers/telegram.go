package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reports-api/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

// SendTelegramNotificationHandler รับ POST แล้วส่งข้อความไป Telegram
func SendTelegramNotificationHandler(c *fiber.Ctx) error {
	var req models.TelegramRequest
	if err := c.BodyParser(&req); err != nil || req.Reportmessage == "" {
		return c.Status(400).JSON(models.TelegramResponse{Success: false, Message: "Invalid request body"})
	}

	// สร้างข้อความรวมข้อมูล
	msg := ""
	if req.BranchName != "" {
		msg += "สาขา: " + req.BranchName + "\n"
	}
	if req.DepartmentName != "" {
		msg += "แผนก: " + req.DepartmentName + "\n"
	}
	if req.Program != "" {
		msg += "โปรแกรม: " + req.Program + "\n"
	}
	if req.Reportmessage != "" {
		msg += "รายงานปัญหา: " + req.Reportmessage + "\n"
	}
	if req.URL != "" {
		msg += "[ดูรายละเอียดเพิ่มเติม](" + req.URL + ")\n"
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")

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
		return c.Status(500).JSON(models.TelegramResponse{Success: false, Message: "Failed to send telegram message: " + err.Error()})
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return c.Status(502).JSON(models.TelegramResponse{Success: false, Message: "Telegram API error: " + string(body)})
	}

	log.Printf("Telegram message sent successfully")
	return c.JSON(models.TelegramResponse{Success: true, Message: "Notification sent successfully"})
}

func SendTelegram(req models.TaskRequest) (int, error) {
	// botToken := os.Getenv("BOT_TOKEN")
	// chatIDStr := os.Getenv("CHAT_ID")

	// chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)

	botToken := "7852676725:AAHnEZclQ57Wo-klSyhZSmbghCU5w0TXgCk"
	chatID := int64(-1002816577414)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, err
	}

	// สร้างข้อความตามสถานะ
	var statusIcon, statusText, headerColor string
	switch req.Status {
	case 0:
		statusIcon = "🔴"
		statusText = "รอดำเนินการ"
		headerColor = "🚨 *แจ้งเตือนปัญหาระบบ* 🚨"
	case 1:
		statusIcon = "✅"
		statusText = "เสร็จสิ้น"
		headerColor = "✅ *งานเสร็จสิ้นแล้ว* ✅"
	}

	msg := headerColor + "\n"
	msg += "━━━━━━━━━━━━━━━━━━━━━━━━\n"

	if req.BranchName != "" {
		msg += "🏢 *สาขา:* `" + req.BranchName + "`\n"
	}
	if req.DepartmentName != "" {
		msg += "🏛️ *แผนก:* `" + req.DepartmentName + "`\n"
	}
	if req.PhoneNumber > 0 {
		msg += fmt.Sprintf("📞 *เบอร์โทร:* `%d`\n", req.PhoneNumber)
	}
	if req.ProgramName != "" {
		msg += "💻 *โปรแกรม:* `" + req.ProgramName + "`\n"
	}
	if req.CreatedAt != "" {
		msg += "📅 *วันที่:* `" + req.CreatedAt + "`\n"
	}

	msg += "\n" + statusIcon + " *สถานะ:* `" + statusText + "`\n"
	msg += "━━━━━━━━━━━━━━━━━━━━━━━━\n"
	msg += "📝 *รายละเอียดปัญหา:*\n"
	msg += "```\n" + req.Text + "\n```"

	if req.Url != "" {
		msg += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = "Markdown"
	log.Printf("%s", message)
	sentMsg, err := bot.Send(message)
	if err != nil {
		return 0, err
	}

	log.Printf("Telegram message sent successfully with ID: %d", sentMsg.MessageID)
	return sentMsg.MessageID, nil
}

func UpdateTelegram(req models.TaskRequest) (int, error) {
	botToken := "7852676725:AAHnEZclQ57Wo-klSyhZSmbghCU5w0TXgCk"
	chatID := int64(-1002816577414)
	messageID := req.MessageID

	// สร้างข้อความตามสถานะ
	var statusIcon, statusText, headerColor string
	switch req.Status {
	case 0:
		statusIcon = "🔴"
		statusText = "รอดำเนินการ"
		headerColor = "🚨 *แจ้งเตือนปัญหาระบบ* 🚨"
	case 1:
		statusIcon = "✅"
		statusText = "เสร็จสิ้น"
		headerColor = "✅ *งานเสร็จสิ้นแล้ว* ✅"
	}

	newMessage := headerColor + "\n"
	newMessage += "━━━━━━━━━━━━━━━━━━━━━━━━\n"

	if req.BranchName != "" {
		newMessage += "🏢 *สาขา:* `" + req.BranchName + "`\n"
	}
	if req.DepartmentName != "" {
		newMessage += "🏛️ *แผนก:* `" + req.DepartmentName + "`\n"
	}
	if req.PhoneNumber > 0 {
		newMessage += fmt.Sprintf("📞 *เบอร์โทร:* `%d`\n", req.PhoneNumber)
	}
	if req.ProgramName != "" {
		newMessage += "💻 *โปรแกรม:* `" + req.ProgramName + "`\n"
	}
	if req.CreatedAt != "" {
		newMessage += "📅 *วันที่:* `" + req.CreatedAt + "`\n"
	}

	newMessage += "\n" + statusIcon + " *สถานะ:* `" + statusText + "`\n"
	newMessage += "━━━━━━━━━━━━━━━━━━━━━━━━\n"
	newMessage += "📝 *รายละเอียดปัญหา:*\n"
	newMessage += "```\n" + req.Text + "\n```"

	if req.Url != "" {
		newMessage += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newMessage)
	editMsg.ParseMode = "Markdown"
	_, err = bot.Send(editMsg)
	if err != nil {
		log.Printf("Error editing message: %v", err)
		return 0, err
	}
	log.Printf("Message ID %d edited successfully!", messageID)
	return messageID, nil
}
