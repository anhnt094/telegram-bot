package cmd

import (
	"bot/common"
	"bot/component/k8s"
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

	bot, err := common.CreateNewBot(conf.TelegramBotToken)
	if err != nil {
		return err
	}

	if err := common.ListenGroupUpdates(bot, conf.TelegramGroupId, clientSet); err != nil {
		return err
	}

	return nil
}
