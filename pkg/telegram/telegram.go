package client

import (
	tele "gopkg.in/telebot.v3"
	"log"
)

type TelegramAPI interface {
	SendMessage(userID uint64, message string, markup *tele.ReplyMarkup) (*tele.Message, error)
}

func NewTelegramApi(token string) TelegramAPI {
	preference := tele.Settings{
		Token:       token,
		ParseMode:   tele.ModeHTML,
		Synchronous: false,
	}
	bot, err := tele.NewBot(preference)
	if err != nil {
		log.Fatal(err)
	}

	return &Telegram{bot: bot}
}

type Telegram struct {
	bot *tele.Bot
}

func (tg *Telegram) SendMessage(userID uint64, message string, markup *tele.ReplyMarkup) (*tele.Message, error) {
	return tg.bot.Send(&tele.User{ID: int64(userID)}, message, markup)

}
