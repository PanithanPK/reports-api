package config

import "os"

type Config struct {
	EndPoint        string
	AccessKey       string
	SecretAccessKey string
	BucketName      string
	Environment     string
	ChatID          string
	BotToken        string
}

var AppConfig *Config

func InitConfig() {
	AppConfig = &Config{
		EndPoint:        os.Getenv("End_POINT"),
		AccessKey:       os.Getenv("ACCESS_KEY"),
		SecretAccessKey: os.Getenv("SECRET_ACCESSKEY"),
		BucketName:      os.Getenv("BUCKET_NAME"),
		Environment:     os.Getenv("env"),
		BotToken:        os.Getenv("BOT_TOKEN"),
		ChatID:          os.Getenv("CHAT_ID"),
	}
}
