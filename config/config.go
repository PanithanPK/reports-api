package config

import (
	"os"
	"reports-api/models"
)

var AppConfig *models.Config

func InitConfig() {
	AppConfig = &models.Config{
		EndPoint:        os.Getenv("End_POINT"),
		AccessKey:       os.Getenv("ACCESS_KEY"),
		SecretAccessKey: os.Getenv("SECRET_ACCESSKEY"),
		BucketName:      os.Getenv("BUCKET_NAME"),
		Environment:     os.Getenv("env"),
		BotToken:        os.Getenv("BOT_TOKEN"),
		ChatID:          os.Getenv("CHAT_ID"),
	}
}
