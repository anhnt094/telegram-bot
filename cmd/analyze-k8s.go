package cmd

import (
	"bot/common"
	"bot/component/k8s"
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

	bot, err := common.CreateNewBot(conf.TelegramBotToken)
	if err != nil {
		return err
	}

	result, err := k8s.Analyze(clientSet)
	if err != nil {
		return err
	}

	if err := common.SendMessageToGroup(bot, conf.TelegramGroupId, result); err != nil {
		return err
	}

	return nil
}
