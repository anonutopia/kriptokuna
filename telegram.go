package main

import (
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v2"
)

func initTelegramBot() *telebot.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:     conf.TelegramAPIKey,
		Poller:    &telebot.LongPoller{Timeout: TelPollerTimeout * time.Second},
		Verbose:   conf.Debug,
		ParseMode: telebot.ModeHTML,
	})

	if err != nil {
		log.Fatal(err)
	}

	return b
}

func logTelegram(message string) {
	group := &telebot.Chat{ID: TelAnonOps}
	if _, err := bot.Send(group, message); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
}

func messageTelegram(message string, groupId int) {
	var group *telebot.Chat
	if conf.Dev {
		group = &telebot.Chat{ID: TelAnonOps}
	} else {
		group = &telebot.Chat{ID: int64(groupId)}
	}
	if _, err := bot.Send(group, message); err != nil {
		log.Println(err)
		logTelegram(err.Error())
	}
}
