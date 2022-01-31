package cmd

import (
	"bot/component/binance"
	"bot/component/telegram"
	"bot/config"
	"fmt"
)

func getUsdtPrice() error {
	conf, err := config.GetConfigs()
	if err != nil {
		return err
	}

	bot, err := telegram.CreateNewBot(conf.TelegramBotToken)
	if err != nil {
		return err
	}

	price, err := binance.GetP2PUsdtHighestPrice()
	if err != nil {
		return err
	}

	result := fmt.Sprintf("%.0f", price)

	if err := telegram.SendMessageToGroup(bot, conf.TelegramGroupId, result); err != nil {
		return err
	}

	return nil
}
