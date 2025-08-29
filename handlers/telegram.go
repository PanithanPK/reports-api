package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reports-api/models"
	"strconv"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTelegram(req models.TaskRequest, photoURL ...string) (int, string, error) {
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
		return 0, "", err
	}
	var Program string
	if req.SystemID > 0 {
		Program = req.ProgramName
	} else {
		Program = req.IssueElse
	}

	bot.Debug = true
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
	msg += "━━━━━━━━━━━━━━\n"

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
	if Program != "" {
		msg += "💻 *โปรแกรม:* `" + Program + "`\n"
	}
	if req.ReportedBy != "" {
		msg += "👤 *ผู้แจ้ง:* `" + req.ReportedBy + "`\n"
	}
	msg += "📅 *วันที่แจ้งปัญหา:* `" + req.CreatedAt + "`\n"
	msg += "━━━━━━━━━━━━━━"
	msg += "\n" + statusIcon + " *สถานะ:* `" + statusText + "`\n"
	if req.Status == 1 {
		msg += "📅 *วันที่แก้ไขเสร็จ:* `" + req.UpdatedAt + "`\n"
	}

	msg += "━━━━━━━━━━━━━━\n"
	msg += "📝 *รายละเอียดปัญหา:*\n"
	msg += "```\n" + req.Text + "\n```"
	msg += "━━━━━━━━━━━━━━"
	// แสดงลิงก์ดูรูปรายงานปัญหา
	if len(photoURL) > 0 {
		for i, url := range photoURL {
			if url != "" {
				msg += fmt.Sprintf("\n🖼️ [ดูรูปรายงานปัญหา %d](%s)", i+1, url)
			}
		}
	}
	msg += "━━━━━━━━━━━━━━"
	if req.Url != "" {
		msg += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}

	var sentMsg tgbotapi.Message
	if len(photoURL) > 0 && photoURL[0] != "" {
		resp, err := http.Get(photoURL[0])
		if err != nil {
			log.Printf("Error fetching photo: %v", err)
			message := tgbotapi.NewMessage(chatID, msg)
			message.ParseMode = "Markdown"
			sentMsg, err = bot.Send(message)
			if err != nil {
				return 0, "", err
			}
		} else {
			defer resp.Body.Close()
			var buf bytes.Buffer
			_, err = io.Copy(&buf, resp.Body)
			if err != nil {
				message := tgbotapi.NewMessage(chatID, msg)
				message.ParseMode = "Markdown"
				sentMsg, err = bot.Send(message)
				if err != nil {
					return 0, "", err
				}
			} else {
				log.Printf("URL Images: %s", photoURL[0])
				photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileReader{
					Name:   photoURL[0],
					Reader: &buf,
				})
				photoMsg.Caption = msg
				photoMsg.ParseMode = "Markdown"
				sentMsg, err = bot.Send(photoMsg)
				if err != nil {
					log.Printf("❌ ส่งภาพไม่สำเร็จ ส่งเป็นข้อความแทน: %v", err)
					message := tgbotapi.NewMessage(chatID, msg)
					message.ParseMode = "Markdown"
					sentMsg, err = bot.Send(message)
					if err != nil {
						return 0, "", err
					}
				}
			}
		}
	} else {
		message := tgbotapi.NewMessage(chatID, msg)
		message.ParseMode = "Markdown"
		sentMsg, err = bot.Send(message)
		if err != nil {
			return 0, "", err
		}
	}

	log.Printf("Telegram message sent successfully with ID: %d", sentMsg.MessageID)
	return sentMsg.MessageID, sentMsg.From.UserName, nil
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

	var Program string
	if req.SystemID > 0 {
		Program = req.ProgramName
	} else {
		Program = req.IssueElse
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

	newMessage := headerColor + "\n"
	newMessage += "━━━━━━━━━━━━━━\n"

	if req.Ticket != "" {
		newMessage += "🎫 *Ticket No:* `" + req.Ticket + "`\n"
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
	if Program != "" {
		newMessage += "💻 *โปรแกรม:* `" + Program + "`\n"
	}
	if req.ReportedBy != "" {
		newMessage += "\n👤 *ผู้แจ้ง:* `" + req.ReportedBy + "`\n"
	}

	newMessage += "📅 *วันที่แจ้งปัญหา:* `" + req.CreatedAt + "`\n"
	newMessage += "━━━━━━━━━━━━━━"
	if req.Assignto != "" {
		newMessage += "👤 *ผู้รับผิดชอบ:* `" + req.Assignto + "`"
	}
	newMessage += "\n" + statusIcon + " *สถานะ:* `" + statusText + "`\n"
	if req.Status == 1 {
		newMessage += "📅 *วันที่แก้ไขเสร็จ:* `" + req.UpdatedAt + "`\n"
	}

	newMessage += "━━━━━━━━━━━━━━\n"
	newMessage += "📝 *รายละเอียดปัญหา:*\n"
	newMessage += "```\n" + req.Text + "\n```"
	newMessage += "━━━━━━━━━━━━━━"
	if len(photoURL) > 0 {
		for i, url := range photoURL {
			if url != "" {
				newMessage += fmt.Sprintf("\n🖼️ [ดูรูปรายงานปัญหา %d](%s)", i+1, url)
			}
		}
	}
	newMessage += "━━━━━━━━━━━━━━"
	if req.Url != "" {
		newMessage += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	// Edit photo caption if photoURL is provided, otherwise edit text message
	if len(photoURL) > 0 && photoURL[0] != "" {
		editMsg := tgbotapi.NewEditMessageCaption(chatID, messageID, newMessage)
		editMsg.ParseMode = "Markdown"
		_, err = bot.Send(editMsg)
	} else {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newMessage)
		editMsg.ParseMode = "Markdown"
		_, err = bot.Send(editMsg)
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

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err = bot.Send(deleteMsg)
	if err != nil {
		log.Printf("Error deleting message: %v", err)
		return false, err
	}

	log.Printf("Message ID %d deleted successfully!", messageID)
	return true, nil
}
