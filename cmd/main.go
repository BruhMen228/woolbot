package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"github.com/BruhMen228/woolbot"
	"github.com/BruhMen228/woolbot/internal/handlers"
	botmth "github.com/BruhMen228/woolbot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := loadEnv()
	if err != nil {
		logger.Error("ошибка загрузки .env", slog.String("ошибка", err.Error()))
	}
	startBot(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, nil)
}

func loadEnv() (error) {
	err := godotenv.Load()
	if err != nil {
		return errors.New("error loading .env file: " + err.Error())
	}
	return nil
}

func startBot(logger *slog.Logger) {
	token := os.Getenv("TOKEN")
	if token == "" {
		logger.Error("ошибка получения токена бота", slog.String("ошибка", "пустой токен"))
		return
	}

	bot, updates, err := woolbot.InitBot(token, false, 60)
	if err != nil {
		logger.Error("ошибка инициализации бота", slog.String("ошибка", err.Error()))
		return
	}

	logger.Info("Бот запущен", slog.String("имя", bot.Self.UserName))

	for update := range updates {
		if update.Message == nil {
			continue
		}		
		
		go func(update tgbotapi.Update) {
			logger.Info("Получено сообщение", 
				slog.String("сообщение", update.Message.Text), 
				slog.String("от кого", update.Message.From.UserName),
			)

			err := handlers.HandleCommand(bot, update)

			if err != nil {
				logger.Error("сообщение не было отправлено", slog.String("ошибка", err.Error()))
				if err.Error() == "context deadline exceeded (Client.Timeout or context cancellation while reading body)" {
					
					text := "Время ожидания ответа превышено. Пожалуйста, попробуйте ещё раз."

					err := botmth.SendMessage(bot, update.Message.Chat.ID, text, update.Message.MessageID)

					if err != nil {
						logger.Error("сообщение бота не было отправлено", slog.String("ошибка", err.Error()))
					}
					
				}
			}
		}(update)
	}
}
