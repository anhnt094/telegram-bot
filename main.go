package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

func sendMessageToGroup(bot *tgbotapi.BotAPI, groupId int64, text string) error {
	msg := tgbotapi.NewMessage(groupId, text)
	msg.ParseMode = "html"
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic("Cannot get the token.")
	}
	bot.Debug = true

	testMsg := `<b>bold</b>, <strong>bold</strong>
	<i>italic</i>, <em>italic</em>
	<a href="http://www.example.com/">inline URL</a>
	<code>inline fixed-width code</code>
	<pre>pre-formatted fixed-width code block</pre>`

	if err := sendMessageToGroup(bot, -709270228, testMsg); err != nil {
		log.Panic(err)
	}

	//updateConfig := tgbotapi.NewUpdate(0)
	//updateConfig.Timeout = 30
	//updates := bot.GetUpdatesChan(updateConfig)
	//for update := range updates {
	//	if update.Message == nil { // ignore any non-Message updates
	//		continue
	//	}
	//
	//	// Create a new MessageConfig. We don't have text yet,
	//	// so we leave it empty.
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	//
	//	if !update.Message.IsCommand() {
	//		msg.Text = "Run /help for more information."
	//	}
	//
	//	// Extract the command from the Message.
	//	switch update.Message.Command() {
	//	case "help":
	//		msg.Text = "I understand /sayhi, /status and /whoismygirlfriend."
	//	case "sayhi":
	//		msg.Text = "Hi :)"
	//	case "status":
	//		msg.Text = "I'm ok."
	//	case "whoismygirlfriend":
	//		msg.Text = "Em Thảo yêu dấu chứ ai. Giả vờ hỏi hả mài !! 🥶"
	//	default:
	//		msg.Text = "I don't know that command. Run /help for more information."
	//	}
	//
	//	if _, err := bot.Send(msg); err != nil {
	//		log.Fatalln(err)
	//	}
	//}
}
