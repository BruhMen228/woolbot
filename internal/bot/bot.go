package bot

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string, replyTo int) error {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeMarkdown

	if replyTo != -1 {
		msg.ReplyToMessageID = replyTo
	}

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func RequestAction(bot *tgbotapi.BotAPI, chatID int64, act string, ctx context.Context) {
	action := tgbotapi.NewChatAction(chatID, act)

	LOOP:
		for {
			// if _, err := bot.Request(action); err != nil {
			//	return err
			//}
			bot.Request(action)
			time.Sleep(4 * time.Second)
			select {
			case <- ctx.Done():
				break LOOP
			default:
				continue
			}
		}

	//return nil
}
