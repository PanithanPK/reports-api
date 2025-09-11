package common

import (
	"fmt"
	"reports-api/models"
)

func FormatSolutionMessage(req models.ResolutionReq, photoURLs ...string) string {
	replyText := "🔧 *วิธีการแก้ไข* 🔧\n"
	replyText += "━━━━━━━━━━━━━━\n"
	replyText += "🎫 *Ticket No:* " + req.TicketNo + "\n"

	if req.TelegramUser == "" {
		replyText += "👤 *ผู้รับผิดชอบ:* " + req.Assignto + "\n"
	} else {
		replyText += "👤 *ผู้รับผิดชอบ:* " + req.Assignto + "" + req.TelegramUser + "\n"
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
