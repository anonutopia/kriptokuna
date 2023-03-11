package main

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"runtime"
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

func getCallerInfo() (info string) {

	// pc, file, lineNo, ok := runtime.Caller(2)
	_, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	// funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	return fmt.Sprintf("%s:%d: ", fileName, lineNo)
}

func logTelegram(message string) {
	message = "anote-mobile:" + getCallerInfo() + url.PathEscape(url.QueryEscape(message))

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
