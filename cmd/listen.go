package cmd

import (
	"bot/component/k8s"
	"bot/component/telegram"
	"bot/config"
)

func listen() error {
	conf, err := config.GetConfigs()
	if err != nil {
		return err
	}

	clientSet, err := k8s.Authenticate(conf.KubeConfig)
	if err != nil {
		return err
	}

	bot, err := telegram.CreateNewBot(conf.TelegramBotToken)
	if err != nil {
		return err
	}

	if err := telegram.ListenGroupUpdates(bot, conf.TelegramGroupId, clientSet); err != nil {
		return err
	}

	return nil
}
