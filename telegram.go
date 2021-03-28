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
	var group *telebot.Chat
	if conf.Dev {
		group = &telebot.Chat{ID: TelAnonOps}
	} else {
		group = &telebot.Chat{ID: TelAnonTeam}
	}
	if _, err := bot.Send(group, message); err != nil {
		log.Println(err)
	}
}
