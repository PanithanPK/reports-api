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
	"strings"

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

	bot.Debug = false
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
	// แสดงลิงก์ดูรูปรายงานปัญหา
	if len(photoURL) > 0 {
		msg += "━━━━━━━━━━━━━━"
		for i, url := range photoURL {
			if url != "" {
				msg += fmt.Sprintf("\n🖼️ [ดูรูปรายงานปัญหา %d](%s)", i+1, url)
			}
		}
	}
	msg += "\n━━━━━━━━━━━━━━"
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
	// Helper function to escape Markdown characters
	escapeMarkdown := func(text string) string {
		// Characters that need to be escaped in Telegram Markdown
		replacer := strings.NewReplacer(
			"_", "\\_",
			"*", "\\*",
			"[", "\\[",
			"]", "\\]",
			"(", "\\(",
			")", "\\)",
			"~", "\\~",
			"`", "\\`",
			">", "\\>",
			"#", "\\#",
			"+", "\\+",
			"-", "\\-",
			"=", "\\=",
			"|", "\\|",
			"{", "\\{",
			"}", "\\}",
			".", "\\.",
			"!", "\\!",
		)
		return replacer.Replace(text)
	}

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
		newMessage += "🎫 *Ticket No:* " + req.Ticket + "\n"
	}
	if req.BranchName != "" {
		newMessage += "🏢 *สาขา:* " + req.BranchName + "\n"
	}
	if req.DepartmentName != "" {
		newMessage += "🏛️ *แผนก:* " + req.DepartmentName + "\n"
	}
	if req.PhoneNumber > 0 {
		newMessage += fmt.Sprintf("📞 *เบอร์โทร:* %d\n", req.PhoneNumber)
	}
	if Program != "" {
		newMessage += "💻 *โปรแกรม:* " + Program + "\n"
	}
	if req.ReportedBy != "" {
		newMessage += "\n👤 *ผู้แจ้ง:* " + req.ReportedBy + "\n"
	}

	newMessage += "📅 *วันที่แจ้งปัญหา:* " + req.CreatedAt + "\n"
	newMessage += "━━━━━━━━━━━━━━"
	if req.Assignto != "" {
		if req.TelegramUser != "" {
			// ใช้ @ เพื่อแท็กผู้ใช้ Telegram
			telegramTag := req.TelegramUser
			if !strings.HasPrefix(telegramTag, "@") {
				telegramTag = "@" + telegramTag
			}
			// Escape underscore in telegram username for Markdown
			telegramTag = strings.ReplaceAll(telegramTag, "_", "\\_")
			newMessage += "\n👤 *ผู้รับผิดชอบ:* " + escapeMarkdown(req.Assignto) + " " + telegramTag
		} else {
			newMessage += "\n👤 *ผู้รับผิดชอบ:* " + escapeMarkdown(req.Assignto)
		}
	}
	newMessage += "\n" + statusIcon + " *สถานะ:* " + escapeMarkdown(statusText) + "\n"
	if req.Status == 1 {
		newMessage += "📅 *วันที่แก้ไขเสร็จ:* " + req.UpdatedAt + "\n"
	}

	newMessage += "━━━━━━━━━━━━━━\n"
	newMessage += "📝 *รายละเอียดปัญหา:*\n"
	newMessage += "```\n" + req.Text + "\n```"

	if len(photoURL) > 0 {
		newMessage += "━━━━━━━━━━━━━━"
		for i, url := range photoURL {
			if url != "" {
				newMessage += fmt.Sprintf("\n🖼️ [ดูรูปรายงานปัญหา %d](%s)", i+1, url)
			}
		}
	}
	newMessage += "\n━━━━━━━━━━━━━━"
	if req.Url != "" {
		newMessage += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Failed to create Telegram bot:", err)
	}
	bot.Debug = false

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

	// ส่งการแจ้งเตือนเฉพาะเมื่อมีการเปลี่ยนผู้รับผิดชอบ
	var notificationID int
	if req.TelegramUser != "" && req.PreviousAssignto != req.Assignto {
		telegramTag := req.TelegramUser
		if !strings.HasPrefix(telegramTag, "@") {
			telegramTag = "@" + telegramTag
		}

		var notificationMsg string
		switch req.Status {
		case 0:
			notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* `%s`\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", escapeMarkdown(telegramTag), req.Ticket, req.Url)
		case 1:

		}

		if notificationMsg != "" {
			notifyMsg := tgbotapi.NewMessage(chatID, notificationMsg)
			notifyMsg.ParseMode = "Markdown"
			notifyMsg.ReplyToMessageID = messageID
			notificationResp, err := bot.Send(notifyMsg)
			if err != nil {
				log.Printf("Warning: Failed to send notification: %v", err)
			} else {
				notificationID = notificationResp.MessageID
			}
		}
	}

	log.Printf("Message ID %d edited successfully!", messageID)
	if notificationID > 0 {
		return notificationID, nil
	}
	return notificationID, nil
}

func DeleteTelegram(messageID int) (bool, error) {
	if messageID <= 0 {
		return false, nil
	}

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Printf("Invalid CHAT_ID format: %v", err)
		return false, err
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to create bot: %v", err)
		return false, err
	}

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	resp, err := bot.Request(deleteMsg)
	if err != nil {
		log.Printf("Cannot delete message ID %d: %v", messageID, err)
		return false, nil // Return nil error to prevent cascade failures
	}
	if !resp.Ok {
		log.Printf("Delete message failed for ID %d: %s", messageID, resp.Description)
		return false, nil
	}

	log.Printf("Message ID %d deleted successfully!", messageID)
	return true, nil
}

func formatSolutionMessage(req models.ResolutionReq, photoURLs ...string) string {
	replyText := "🔧 *วิธีการแก้ไข* 🔧\n"
	replyText += "━━━━━━━━━━━━━━\n"
	replyText += "🎫 *Ticket No:* `" + req.TicketNo + "`\n"
	replyText += "👤 *ผู้รับผิดชอบ:* `" + req.Assignto + "`\n"
	replyText += "📅 *วันที่แจ้ง:* `" + req.CreatedAt + "`\n"
	replyText += "📅 *วันที่แก้ไข:* `" + req.ResolvedAt + "`\n"
	replyText += "━━━━━━━━━━━━━━\n"

	replyText += "📝 *รายละเอียดการแก้ไข:*\n"
	replyText += "```\n" + req.Solution + "\n```"

	// Add photo links if available
	if len(photoURLs) > 0 {
		replyText += "\n━━━━━━━━━━━━━━"
		for i := 0; i < len(photoURLs); i++ {
			if photoURLs[i] != "" {
				replyText += fmt.Sprintf("\n🖼️ [ดูรูปการแก้ไข %d](%s)", i+1, photoURLs[i])
			}
		}
	}
	replyText += "\n━━━━━━━━━━━━━━"
	replyText += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")"

	return replyText
}

func replyToSpecificMessage(req models.ResolutionReq, photoURLs ...string) (int, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Format solution message
	replyText := formatSolutionMessage(req, photoURLs...)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, err
	}

	var sentMsg tgbotapi.Message

	// ส่งรูปแรกพร้อมข้อความ ถ้ามีรูป
	if len(photoURLs) > 0 && photoURLs[0] != "" {
		resp, err := http.Get(photoURLs[0])
		if err != nil {
			log.Printf("Error fetching photo: %v", err)
			// ส่งเป็นข้อความแทน
			message := tgbotapi.NewMessage(chatID, replyText)
			message.ParseMode = "Markdown"
			message.ReplyToMessageID = req.MessageID
			sentMsg, err = bot.Send(message)
			if err != nil {
				return 0, err
			}
		} else {
			defer resp.Body.Close()
			var buf bytes.Buffer
			_, err = io.Copy(&buf, resp.Body)
			if err != nil {
				// ส่งเป็นข้อความแทน
				message := tgbotapi.NewMessage(chatID, replyText)
				message.ParseMode = "Markdown"
				message.ReplyToMessageID = req.MessageID
				sentMsg, err = bot.Send(message)
				if err != nil {
					return 0, err
				}
			} else {
				// ส่งรูปพร้อม caption
				photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileReader{
					Name:   photoURLs[0],
					Reader: &buf,
				})
				photoMsg.Caption = replyText
				photoMsg.ParseMode = "Markdown"
				photoMsg.ReplyToMessageID = req.MessageID
				sentMsg, err = bot.Send(photoMsg)
				if err != nil {
					log.Printf("❌ ส่งภาพไม่สำเร็จ ส่งเป็นข้อความแทน: %v", err)
					message := tgbotapi.NewMessage(chatID, replyText)
					message.ParseMode = "Markdown"
					message.ReplyToMessageID = req.MessageID
					sentMsg, err = bot.Send(message)
					if err != nil {
						return 0, err
					}
				}
			}
		}
	} else {
		// ส่งเฉพาะข้อความ
		message := tgbotapi.NewMessage(chatID, replyText)
		message.ParseMode = "Markdown"
		message.ReplyToMessageID = req.MessageID
		sentMsg, err = bot.Send(message)
		if err != nil {
			return 0, err
		}
	}

	log.Printf("Solution reply sent successfully with ID: %d", sentMsg.MessageID)
	return sentMsg.MessageID, nil
}

func UpdatereplyToSpecificMessage(messageID int, req models.ResolutionReq, photoURLs ...string) (int, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	chatIDStr := os.Getenv("CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Format solution message
	replyText := formatSolutionMessage(req, photoURLs...)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, err
	}

	// ถ้ามีรูป ให้ edit caption ของรูปเดิม
	if len(photoURLs) > 0 && photoURLs[0] != "" {
		editMsg := tgbotapi.NewEditMessageCaption(chatID, messageID, replyText)
		editMsg.ParseMode = "Markdown"
		_, err = bot.Send(editMsg)
		if err != nil {
			log.Printf("Error editing photo caption: %v", err)
			return 0, err
		}
	} else {
		// ถ้าไม่มีรูป ให้ edit ข้อความเดิม
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, replyText)
		editMsg.ParseMode = "Markdown"
		_, err = bot.Send(editMsg)
		if err != nil {
			log.Printf("Error editing message text: %v", err)
			return 0, err
		}
	}
	log.Printf("Solution edit sent successfully for message ID: %d", messageID)
	return messageID, nil
}
