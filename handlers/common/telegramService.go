package common

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"reports-api/config"
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

func FormatSolutionMessage(req models.ResolutionReq, photoURLs ...string) string {
	replyText := "🔧 *การแก้ไขปัญหา* 🔧\n"
	replyText += "━━━━━━━━━━━━━━\n"
	replyText += "🎫 *Ticket No:* [" + req.TicketNo + "](" + req.Url + ")\n"

	if req.TelegramUser != "" {
		telegramTag := req.TelegramUser
		if !strings.HasPrefix(telegramTag, "@") {
			telegramTag = "@" + telegramTag
		}
		telegramTag = strings.ReplaceAll(telegramTag, "_", "\\_")
		replyText += "👥 *ผู้รับผิดชอบ:* " + req.Assignto + " " + telegramTag + "\n"
	} else {
		replyText += "👥 *ผู้รับผิดชอบ:* " + req.Assignto + "\n"
	}
	replyText += "📅 *วันที่แจ้ง:* " + req.CreatedAt + "\n"
	replyText += "📅 *วันที่แก้ไข:* " + req.ResolvedAt + "\n"
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

func FormatRepostMessage(req models.TaskRequest, photoURLs ...string) string {

	var Program string
	if req.SystemID > 0 {
		Program = req.ProgramName
	} else {
		Program = req.IssueElse
	}

	// Create message based on status
	var statusIcon, statusText, headerColor string
	switch req.Status {
	case 0:
		statusIcon = "🔴"
		statusText = "รอดำเนินการ"
		headerColor = "🚨 *แจ้งเตือนปัญหาระบบ* 🚨"
	case 1:
		statusIcon = "🔵"
		statusText = "กำลังดำเนินการ"
		headerColor = "🔄 *กำลังดำเนินการแก้ไข* 🔄"
	case 2:
		statusIcon = "✅"
		statusText = "เสร็จสิ้น"
		headerColor = "✅ *งานเสร็จสิ้นแล้ว* ✅"

	}

	newMessage := headerColor + "\n"
	newMessage += "━━━━━━━━━━━━━━\n"

	if req.Ticket != "" {
		newMessage += "🎫 *Ticket No:* [" + req.Ticket + "](" + req.Url + ")\n"
	}
	if req.BranchName != "" {
		newMessage += "🏭 *สาขา:* " + req.BranchName + "\n"
	}
	if req.DepartmentName != "" {
		newMessage += "🏢 *แผนก:* " + req.DepartmentName + "\n"
	}
	if req.PhoneNumber > 0 {
		newMessage += fmt.Sprintf("📠 *เบอร์โทร:* %d\n", req.PhoneNumber)
	} else if req.PhoneNumber == 0 {
		if req.PhoneElse != nil && *req.PhoneElse != "" {
			newMessage += "📠 *เบอร์โทร:* " + *req.PhoneElse + "\n"
		}
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
			// Use @ to tag Telegram user
			telegramTag := req.TelegramUser
			if !strings.HasPrefix(telegramTag, "@") {
				telegramTag = "@" + telegramTag
			}
			// Escape underscore in telegram username for Markdown
			telegramTag = strings.ReplaceAll(telegramTag, "_", "\\_")
			newMessage += "\n👥 *ผู้รับผิดชอบ:* " + EscapeMarkdown(req.Assignto) + " " + telegramTag
		} else {
			newMessage += "\n👥 *ผู้รับผิดชอบ:* " + EscapeMarkdown(req.Assignto)
		}
	}
	newMessage += "\n" + statusIcon + " *สถานะ:* " + EscapeMarkdown(statusText) + "\n"
	if req.Status == 1 {
		newMessage += "📆 *กำลังดำเนินการ:* " + req.UpdatedAt + "\n"
	}
	if req.Status == 2 {
		newMessage += "📅 *วันที่แก้ไขเสร็จ:* " + req.ResolvedAt + "\n"
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
	// Create message based on status
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

			// Process downloaded image for Telegram optimization
			imageConfig := DefaultImageConfig()
			processedImage, _, err := ProcessImageFromReader(resp.Body, photoURL[0], imageConfig)
			if err != nil {
				log.Printf("Error processing downloaded image: %v, sending as text", err)
				message := tgbotapi.NewMessage(chatID, msg)
				message.ParseMode = "Markdown"
				sentMsg, err = bot.Send(message)
				if err != nil {
					return 0, "", err
				}
			} else {
				log.Printf("URL Images processed: %s, Size: %d bytes", photoURL[0], processedImage.Len())
				photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileReader{
					Name:   photoURL[0],
					Reader: bytes.NewReader(processedImage.Bytes()),
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

	return sentMsg.MessageID, sentMsg.From.UserName, nil
}

func UpdateTelegram(req models.TaskRequest, photoURL ...string) (int, error) {
	// Helper function to escape Markdown characters
	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Printf("UpdateTelegram - Invalid CHAT_ID format: %v", err)
		log.Fatal("Invalid CHAT_ID format:", err)
	}
	messageID := req.MessageID

	newMessage := FormatRepostMessage(req, photoURL...)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Failed to create Telegram bot:", err)
	}
	bot.Debug = false

	// Edit photo caption if photoURL is provided, otherwise edit text message
	if len(photoURL) > 0 && photoURL[0] != "" {
		editMsg := tgbotapi.NewEditMessageCaption(chatID, messageID, newMessage)
		editMsg.ParseMode = "Markdown"
		_, err := bot.Send(editMsg)
		if err != nil {
			log.Printf("UpdateTelegram - Error editing message caption: %v", err)
			return 0, err
		}
	} else {
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, newMessage)
		editMsg.ParseMode = "Markdown"
		_, err := bot.Send(editMsg)
		if err != nil {
			log.Printf("UpdateTelegram - Error editing message text: %v", err)
			return 0, err
		}
	}

	// Send notification only when assignee changes
	var notificationID int

	if req.TelegramUser != "" && req.PreviousAssignto != req.Assignto {
		telegramTag := req.TelegramUser
		if !strings.HasPrefix(telegramTag, "@") {
			telegramTag = "@" + telegramTag
		}

		var notificationMsg string
		switch req.Status {
		case 0:
			notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* [%s](%s)\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url, req.Url)
		case 1:
			notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* [%s](%s)\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url, req.Url)
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
				log.Printf("UpdateTelegram - Notification sent with ID: %d", notificationID)
			}
		}
	}

	log.Printf("Message ID %d edited successfully!", messageID)
	if notificationID > 0 {
		return notificationID, nil
	}
	return 0, nil // Return 0 instead of notificationID for clarity
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
		notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* `%s`\n━━━━━━━━━━━━━━\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url)
	case 1:
		notificationMsg = fmt.Sprintf("🔔 *การแจ้งเตือนมอบหมายงาน* 🔔\n━━━━━━━━━━━━━━\n👋 %s\n📋 คุณได้รับมอบหมายงานใหม่แล้ว\n🎫 *Ticket:* `%s`\n━━━━━━━━━━━━━━\n🔗 [ดูรายละเอียดเพิ่มเติม](%s)\n━━━━━━━━━━━━━━", EscapeMarkdown(telegramTag), req.Ticket, req.Url)
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

func ReplyToSpecificMessage(req models.ResolutionReq, photoURLs ...string) (int, error) {
	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	// Format solution message
	replyText := FormatSolutionMessage(req, photoURLs...)
	log.Printf("TelegramUser: %s", req.TelegramUser)

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

			// Process downloaded image for Telegram optimization
			imageConfig := DefaultImageConfig()
			processedImage, _, err := ProcessImageFromReader(resp.Body, photoURLs[0], imageConfig)
			if err != nil {
				log.Printf("Error processing downloaded image: %v, sending as text", err)
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
					Reader: bytes.NewReader(processedImage.Bytes()),
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
	log.Printf("🔄 Starting UpdatereplyToSpecificMessage for messageID: %d", messageID)

	botToken := config.AppConfig.BotToken
	chatIDStr := config.AppConfig.ChatID
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid chat ID: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return 0, fmt.Errorf("failed to create bot: %v", err)
	}

	// ลองลบ message เก่า (ไม่ถือเป็น error ร้ายแรงถ้าลบไม่ได้)
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	resp, err := bot.Request(deleteMsg)
	if err != nil {
		log.Printf("⚠️ Cannot delete message ID %d: %v (continuing anyway)", messageID, err)
	} else if !resp.Ok {
		log.Printf("⚠️ Delete message failed for ID %d: %s (continuing anyway)", messageID, resp.Description)
	}

	// Format solution message
	replyText := FormatSolutionMessage(req, photoURLs...)
	log.Printf("📝 Formatted message length: %d characters", len(replyText))

	var sentMsg tgbotapi.Message

	// ส่ง message ใหม่
	if len(photoURLs) > 0 && photoURLs[0] != "" {
		sentMsg, err = sendPhotoMessage(bot, chatID, photoURLs[0], replyText, req.MessageID)
	} else {
		sentMsg, err = sendTextMessage(bot, chatID, replyText, req.MessageID)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("✅ Solution message sent successfully with new ID: %d", sentMsg.MessageID)
	return sentMsg.MessageID, nil
}

// Helper functions
func sendPhotoMessage(bot *tgbotapi.BotAPI, chatID int64, photoURL, caption string, replyToMessageID int) (tgbotapi.Message, error) {

	resp, err := http.Get(photoURL)
	if err != nil {
		log.Printf("⚠️ Error fetching photo: %v, sending as text instead", err)
		return sendTextMessage(bot, chatID, caption, replyToMessageID)
	}
	defer resp.Body.Close()

	// Process downloaded image for Telegram optimization
	imageConfig := DefaultImageConfig()
	processedImage, _, err := ProcessImageFromReader(resp.Body, photoURL, imageConfig)
	if err != nil {
		log.Printf("⚠️ Error processing photo: %v, sending as text instead", err)
		return sendTextMessage(bot, chatID, caption, replyToMessageID)
	}

	log.Printf("📸 Photo processed successfully, size: %d bytes", processedImage.Len())

	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileReader{
		Name:   photoURL,
		Reader: bytes.NewReader(processedImage.Bytes()),
	})
	photoMsg.Caption = caption
	photoMsg.ParseMode = "Markdown"
	photoMsg.ReplyToMessageID = replyToMessageID

	sentMsg, err := bot.Send(photoMsg)
	if err != nil {
		log.Printf("⚠️ Failed to send photo: %v, sending as text instead", err)
		return sendTextMessage(bot, chatID, caption, replyToMessageID)
	}

	log.Printf("✅ Photo message sent successfully")
	return sentMsg, nil
}

func sendTextMessage(bot *tgbotapi.BotAPI, chatID int64, text string, replyToMessageID int) (tgbotapi.Message, error) {
	log.Printf("📄 Sending text message")

	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "Markdown"
	message.ReplyToMessageID = replyToMessageID

	sentMsg, err := bot.Send(message)
	if err != nil {
		log.Printf("❌ Failed to send text message: %v", err)
		return sentMsg, err
	}

	log.Printf("✅ Text message sent successfully")
	return sentMsg, nil
}
