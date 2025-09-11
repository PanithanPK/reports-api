package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"reports-api/config"
	"reports-api/handlers/common"
	"reports-api/models"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func EscapeMarkdown(text string) string {
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

func FormatRepostMessage(req models.TaskRequest, photoURLs ...string) string {

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
			newMessage += "\n👤 *ผู้รับผิดชอบ:* " + EscapeMarkdown(req.Assignto) + " " + telegramTag
		} else {
			newMessage += "\n👤 *ผู้รับผิดชอบ:* " + EscapeMarkdown(req.Assignto)
		}
	}
	newMessage += "\n" + statusIcon + " *สถานะ:* " + EscapeMarkdown(statusText) + "\n"
	if req.Status == 1 {
		newMessage += "📅 *วันที่แก้ไขเสร็จ:* " + req.UpdatedAt + "\n"
	}

	newMessage += "━━━━━━━━━━━━━━\n"
	newMessage += "📝 *รายละเอียดปัญหา:*\n"
	newMessage += "```\n" + req.Text + "\n```"

	if len(photoURLs) > 0 {
		newMessage += "━━━━━━━━━━━━━━"
		for i, url := range photoURLs {
			if url != "" {
				newMessage += fmt.Sprintf("\n🖼️ [ดูรูปรายงานปัญหา %d](%s)\n", i+1, url)
			}
		}
	}
	newMessage += "\n━━━━━━━━━━━━━━"
	if req.Url != "" {
		newMessage += "\n🔗 [ดูรายละเอียดเพิ่มเติม](" + req.Url + ")\n"
	}

	return newMessage
}

func SendTelegram(req models.TaskRequest, photoURL ...string) (int, string, error) {
	// botToken := os.Getenv("BOT_TOKEN")
	// chatIDStr := os.Getenv("CHAT_ID")

	// chatID, _ := strconv.ParseInt(chatIDStr, 10, 64)
	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid CHAT_ID format:", err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, "", err
	}

	bot.Debug = false
	// สร้างข้อความตามสถานะ
	msg := FormatRepostMessage(req, photoURL...)

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
	log.Printf("UpdateTelegram called - MessageID: %d, Status: %d, Assignto: %s", req.MessageID, req.Status, req.Assignto)
	log.Printf("UpdateTelegram - PhotoURLs: %v", photoURL)

	// Helper function to escape Markdown characters
	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Printf("UpdateTelegram - Invalid CHAT_ID format: %v", err)
		log.Fatal("Invalid CHAT_ID format:", err)
	}
	messageID := req.MessageID

	log.Printf("UpdateTelegram - Bot config: chatID=%d, messageID=%d", chatID, messageID)

	newMessage := FormatRepostMessage(req, photoURL...)
	log.Printf("UpdateTelegram - Formatted message length: %d characters", len(newMessage))

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("UpdateTelegram - Failed to create Telegram bot: %v", err)
		log.Panic("Failed to create Telegram bot:", err)
	}
	bot.Debug = false

	// Edit photo caption if photoURL is provided, otherwise edit text message
	if len(photoURL) > 0 && photoURL[0] != "" {
		log.Printf("UpdateTelegram - Editing message caption for messageID: %d", messageID)
		editMsg := tgbotapi.NewEditMessageCaption(chatID, messageID, newMessage)
		editMsg.ParseMode = "Markdown"
		resp, err := bot.Send(editMsg)
		if err != nil {
			log.Printf("UpdateTelegram - Error editing message caption: %v", err)
			return 0, err
		}
		log.Printf("UpdateTelegram - Caption edit response: %+v", resp)
	} else {
		log.Printf("UpdateTelegram - Editing message text for messageID: %d", messageID)
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newMessage)
		editMsg.ParseMode = "Markdown"
		_, err := bot.Send(editMsg)
		if err != nil {
			log.Printf("UpdateTelegram - Error editing message text: %v", err)
			return 0, err
		}
	}

	log.Printf("UpdateTelegram - Message edit successful")

	// ส่งการแจ้งเตือนเฉพาะเมื่อมีการเปลี่ยนผู้รับผิดชอบ
	var notificationID int
	log.Printf("UpdateTelegram - Checking notification: TelegramUser='%s', PreviousAssignto='%s', CurrentAssignto='%s'",
		req.TelegramUser, req.PreviousAssignto, req.Assignto)

	if req.TelegramUser != "" && req.PreviousAssignto != req.Assignto {
		log.Printf("UpdateTelegram - Sending notification for assignee change")
		telegramTag := req.TelegramUser
		if !strings.HasPrefix(telegramTag, "@") {
			telegramTag = "@" + telegramTag
		}

		var notificationMsg string
		switch req.Status {
		case 0:
			notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* `%s`\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url)
		case 1:
			// ไม่มี case 1 ใน original code
		}

		if notificationMsg != "" {
			log.Printf("UpdateTelegram - Sending notification message")
			notifyMsg := tgbotapi.NewMessage(chatID, notificationMsg)
			notifyMsg.ParseMode = "Markdown"
			notifyMsg.ReplyToMessageID = messageID
			notificationResp, err := bot.Send(notifyMsg)
			if err != nil {
				log.Printf("Warning: Failed to send notification: %v", err)
			} else {
				notificationID = notificationResp.MessageID
				log.Printf("UpdateTelegram - Notification sent with ID: %d", notificationID)
			}
		} else {
			log.Printf("UpdateTelegram - No notification message for status: %d", req.Status)
		}
	} else {
		log.Printf("UpdateTelegram - No notification needed")
	}

	log.Printf("Message ID %d edited successfully!", messageID)
	if notificationID > 0 {
		return notificationID, nil
	}
	return 0, nil // เปลี่ยนจาก notificationID เป็น 0 เพื่อความชัดเจน
}

func UpdateAssignedtoMsg(messageID int, req models.TaskRequest) (int, error) {

	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, err
	}

	telegramTag := req.TelegramUser
	if telegramTag != "" && !strings.HasPrefix(telegramTag, "@") {
		telegramTag = "@" + telegramTag
	}

	var notificationMsg string
	switch req.Status {
	case 0:
		notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* `%s`\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url)
	case 1:
		notificationMsg = fmt.Sprintf("✅ *งานเสร็จสิ้นแล้ว* ✅\n━━━━━━━━━━━━━━\n👋 %s\n📋 งานที่คุณรับผิดชอบเสร็จสิ้นแล้ว\n🎫 *Ticket:* `%s`\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url)
	}

	if messageID > 0 {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, notificationMsg)
		editMsg.ParseMode = "Markdown"
		_, err = bot.Send(editMsg)
		if err != nil {
			return 0, err
		}
		return messageID, nil
	} else {
		notifyMsg := tgbotapi.NewMessage(chatID, notificationMsg)
		notifyMsg.ParseMode = "Markdown"
		notifyMsg.ReplyToMessageID = req.MessageID
		notificationResp, err := bot.Send(notifyMsg)
		if err != nil {
			return 0, err
		}
		return notificationResp.MessageID, nil
	}
}

func DeleteTelegram(messageID int) (bool, error) {
	if messageID <= 0 {
		return false, nil
	}

	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID

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

func replyToSpecificMessage(req models.ResolutionReq, photoURLs ...string) (int, error) {
	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Format solution message
	replyText := common.FormatSolutionMessage(req, photoURLs...)

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

	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Format solution message
	replyText := common.FormatSolutionMessage(req, photoURLs...)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, err
	}

	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	resp, err := bot.Request(deleteMsg)
	if err != nil {
		log.Printf("Cannot delete message ID %d: %v", messageID, err)
		return 0, nil // Return nil error to prevent cascade failures
	}
	if !resp.Ok {
		log.Printf("Delete message failed for ID %d: %s", messageID, resp.Description)
		return 0, nil
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

	log.Printf("Solution edit sent successfully for message ID: %d", sentMsg.MessageID)
	return sentMsg.MessageID, nil
}
