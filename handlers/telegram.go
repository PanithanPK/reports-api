package handlers

import (
	"fmt"
	"log"
	"os"
	"reports-api/models"
	"strconv"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegram(req models.TaskRequest, photoURL ...string) (int, error) {
	// botToken := os.Getenv("BOT_TOKEN")
	// chatIDStr := os.Getenv("CHAT_ID")

	// chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid CHAT_ID format:", err)
	}

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

	if req.Ticket != "" {
		msg += "🎫 *Ticket No:* `" + req.Ticket + "`\n"
	}
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
	if req.ReportedBy != "" {
		msg += "👤 *ผู้แจ้ง:* `" + req.ReportedBy + "`\n"
	}
	msg += "📅 *วันที่แจ้งปัญหา:* `" + req.CreatedAt + "`\n"
	msg += "━━━━━━━━━━━━━━━━━━━━━━━━"
	msg += "\n" + statusIcon + " *สถานะ:* `" + statusText + "`\n"
	if req.Status == 1 {
		msg += "📅 *วันที่แก้ไขเสร็จ:* `" + req.UpdatedAt + "`\n"
	}

	msg += "━━━━━━━━━━━━━━━━━━━━━━━━\n"
	msg += "📝 *รายละเอียดปัญหา:*\n"
	msg += "```\n" + req.Text + "\n```"

	if req.Url != "" {
		msg += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}

	// Always send text message first
	message := tgbotapi.NewMessage(chatID, msg)
	message.ParseMode = "Markdown"
	sentMsg, err := bot.Send(message)
	if err != nil {
		return 0, err
	}

	// Send photo separately if photoURL is provided
	if len(photoURL) > 0 && photoURL[0] != "" {
		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(photoURL[0]))
		_, _ = bot.Send(photoMsg)
	}

	if err != nil {
		return 0, err
	}

	log.Printf("Telegram message sent successfully with ID: %d", sentMsg.MessageID)
	return sentMsg.MessageID, nil
}

func UpdateTelegram(req models.TaskRequest, photoURL ...string) (int, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid CHAT_ID format:", err)
	}
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

	if req.Ticket != "" {
		newMessage += "\n🎫 *Ticket No:* `" + req.Ticket + "`\n"
	}
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
	if req.ReportedBy != "" {
		newMessage += "👤 *ผู้แจ้ง:* `" + req.ReportedBy + "`\n"
	}

	newMessage += "📅 *วันที่แจ้งปัญหา:* `" + req.CreatedAt + "`\n"
	newMessage += "━━━━━━━━━━━━━━━━━━━━━━━━\n"
	if req.Assignto != "" {
		newMessage += "\n👤 *ผู้รับผิดชอบ:* `" + req.Assignto + "`"
	}
	newMessage += "\n" + statusIcon + " *สถานะ:* `" + statusText + "`\n"
	if req.Status == 1 {
		newMessage += "📅 *วันที่แก้ไขเสร็จ:* `" + req.UpdatedAt + "`\n"
	}

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

	// Try editing text message first, if it fails try editing caption
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newMessage)
	editMsg.ParseMode = "Markdown"
	_, err = bot.Send(editMsg)

	// If editing text fails, try editing caption (for photo messages)
	if err != nil {
		editCaption := tgbotapi.NewEditMessageCaption(chatID, messageID, newMessage)
		editCaption.ParseMode = "Markdown"
		_, err = bot.Send(editCaption)
	}

	if err != nil {
		log.Printf("Error editing message: %v", err)
		return 0, err
	}
	log.Printf("Message ID %d edited successfully!", messageID)
	return messageID, nil
}

func DeleteTelegram(messageID int) (bool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid CHAT_ID format:", err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return false, err
	}

	// Delete main message
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err = bot.Send(deleteMsg)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return false, err
	}

	// Try to delete photo message (next message ID)
	deletePhotoMsg := tgbotapi.NewDeleteMessage(chatID, messageID+1)
	_, _ = bot.Send(deletePhotoMsg) // Ignore error if photo doesn't exist

	log.Printf("Message ID %d deleted successfully!", messageID)
	return true, nil
}
