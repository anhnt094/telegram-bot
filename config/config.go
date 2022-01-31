package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	TelegramBotToken string
	TelegramGroupId  int64
	KubeConfig       string // path to kubeconfig file
}

func GetConfigs() (*Config, error) {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken == "" {
		return nil, errors.New("can't get TELEGRAM_BOT_TOKEN")
	}

	telegramGroupId, err := strconv.Atoi(os.Getenv("TELEGRAM_GROUP_ID"))
	if err != nil {
		return nil, errors.New("can't get TELEGRAM_GROUP_ID")
	}

	kubeconfig := os.Getenv("KUBECONFIG")

	return &Config{
		TelegramBotToken: telegramBotToken,
		TelegramGroupId:  int64(telegramGroupId),
		KubeConfig:       kubeconfig,
	}, nil
}
