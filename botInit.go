package woolbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitBot(token string, debug bool, timeout int) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, nil, err
	}

	bot.Debug = debug

	u := tgbotapi.NewUpdate(0)

	if timeout != 0 {
   		u.Timeout = timeout
	}
	
    updates := bot.GetUpdatesChan(u)

	return bot, updates, nil
}
