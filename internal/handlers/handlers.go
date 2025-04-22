package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/BruhMen228/woolbot/internal/openRouter"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(bot *tgbotapi.BotAPI, update tgbotapi.Update) error

var Commands = map[string]CommandHandler{
	"start": StartHandler,
	"help": HelpHandler,
	"info": InfoHandler,
}

func HandleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	if update.Message != nil {
		if update.Message.IsCommand() {
			if handler, ok := Commands[update.Message.Command()]; ok {
				err := handler(bot, update)
				return err
			}
		}
		if update.Message.Text != "" {
			err := TextHandler(bot, update)
			return err
		}
	}
	return nil
}

func StartHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	action := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	if _, err := bot.Request(action); err != nil {
		return err
	}

	msgFrom := update.Message.From.FirstName

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет %s! Это WoolBot. Чтобы задать вопрос, напиши: \"Шерсть, (твой вопрос).\"", msgFrom))

	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func HelpHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	action := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	if _, err := bot.Request(action); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Доступные команды:\n/start - начать\n/help - помощь\n/info - информация о боте")

	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func InfoHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	action := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	if _, err := bot.Request(action); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "WoolBot - секретное оружие Великих Шерстистых Повторителей по подавлению восстаний\nТакже в WoolBot встроена нейросеть для общения и ответов на вопросы")

	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

type APIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type APIResponseError struct{
	Error struct{
		Message string `json:"message"`
	} `json:"error"`
}


func TextHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

	msgText := update.Message.Text

	msgSplitted:= (strings.Split(msgText, " "))

	if len(msgSplitted) < 2 {
		return errors.New("неправильное сообщение")
	}

	name := strings.ToLower(msgSplitted[0])

	if !(strings.Contains(name, "шерсть") || 
		strings.Contains(name, "woolbot") ||
		strings.Contains(name, "wool") ||
		strings.Contains(name, "шерстьбот") || 
		strings.Contains(name, "шерстистый") ||
		strings.Contains(name, "шерстистыйбот")) {
		return errors.New("неправильное сообщение")
	}

	msgQuery := strings.Join(msgSplitted[1:], " ")

	action := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)

	if _, err := bot.Request(action); err != nil {
		return err
	}

	ctx, err := os.ReadFile("./История клана.txt")
	if err != nil {
		return err
	}
	
	apiKey := os.Getenv("API_KEY2")

	resp, err := openrouter.RequestToOpenRouterAi(apiKey, "deepseek/deepseek-chat-v3-0324:free", "Targon", string(ctx), msgQuery)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var textData APIResponse

	err = json.Unmarshal(respBody, &textData)
	if err != nil {
		return err
	}

	var respAi string

	if textData.Choices != nil {
		respAi = textData.Choices[0].Message.Content
	} else {
		var textData APIResponseError

		err = json.Unmarshal(respBody, &textData)
		if err != nil {
			return err
		}

		respAi = textData.Error.Message
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, respAi)

	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
