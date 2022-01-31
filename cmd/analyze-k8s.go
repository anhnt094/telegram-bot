package cmd

import (
	"bot/component/k8s"
	"bot/component/telegram"
	"bot/config"
)

func analyzeK8s() error {
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

	result, err := k8s.Analyze(clientSet)
	if err != nil {
		return err
	}

	if err := telegram.SendMessageToGroup(bot, conf.TelegramGroupId, result); err != nil {
		return err
	}

	return nil
}
