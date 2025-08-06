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

	botToken := "7852676725:AAHnEZclQ57Wo-klSyhZSmbghCU5w0TXgCk"
	chatID := "-1002816577414"

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

func SendTelegram(req models.TaskRequest) error {
	botToken := "7852676725:AAHnEZclQ57Wo-klSyhZSmbghCU5w0TXgCk"
	chatID := int64(-1002816577414)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return err
	}

	// สร้างข้อความ
	msg := "🚨 *แจ้งเตือนปัญหาระบบ* 🚨\n"

	if req.BranchName != "" {
		msg += "\n🏢 *สาขา* `" + "`\n`" + req.BranchName + "`\n"
	}
	if req.DepartmentName != "" {
		msg += "\n🏢 *แผนก* `" + "`\n`" + req.DepartmentName + "`\n"
	}
	if req.PhoneNumber > 0 {
		msg += fmt.Sprintf("📞 *เบอร์โทร:* `%d`\n", req.PhoneNumber)
	}
	if req.ProgramName != "" {
		msg += "\n💻 *โปรแกรม* `" + "`\n`" + req.ProgramName + "`\n"
	}
	if req.CreatedAt != "" {
		msg += "\n📅 *วันที่* `" + "`\n`" + req.CreatedAt + "`\n"
	}

	msg += "\n⚠️ *รายงานปัญหา:*\n"
	msg += "```\n" + req.Text + "\n```\n"

	if req.Url != "" {
		msg += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = "Markdown"

	_, err = bot.Send(message)
	if err != nil {
		return err
	}

	log.Printf("Telegram message sent successfully")
	return nil
}
