package telegram

import (
	"bot/component/binance"
	"bot/component/k8s"
	"bytes"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"k8s.io/client-go/kubernetes"
	"log"
	"text/template"
)

func CreateNewBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = false

	return bot, err
}

func ListenGroupUpdates(bot *tgbotapi.BotAPI, groupId int64, clientSet *kubernetes.Clientset) error {
	log.Println("listening group updates...")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	var err error
	var respText string

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if update.Message.Chat.ID != groupId {
			continue
		}

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			var output bytes.Buffer

			var tmpl *template.Template
			tmpl = template.Must(template.ParseFiles("templates/telegram-help.txt"))
			if err = tmpl.Execute(&output, nil); err != nil {
				return err
			}

			respText = output.String()
		case "sayhi":
			respText = "Hi :)"
		case "k8s":
			respText, err = k8s.Analyze(clientSet)
			if err != nil {
				return err
			}
		case "usdt":
			respText, err = binance.GetP2PUsdtHighestPriceReport()
			if err != nil {
				return err
			}
		default:
			respText = "I don't know that command. Run /help for more information."
		}

		if err := SendMessageToGroup(bot, groupId, respText); err != nil {
			return err
		}
	}
	return nil
}

func SendMessageToGroup(bot *tgbotapi.BotAPI, groupId int64, text string) error {
	log.Println("sending message to Telegram")
	msg := tgbotapi.NewMessage(groupId, text)
	msg.ParseMode = "html"
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil
}
