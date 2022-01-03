package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /sayhi, /status and /whoismygirlfriend."
		case "sayhi":
			msg.Text = "Hi :)"
		case "status":
			msg.Text = "I'm ok."
		case "whoismygirlfriend":
			msg.Text = "Em Thảo yêu dấu chứ ai. Giả vờ hỏi hả mài !! 🥶"
		default:
			msg.Text = "I don't know that command. Run /help for more information."
		}

		if !update.Message.IsCommand() {
			msg.Text = "Run /help for more information."
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatalln(err)
		}

	}
}
