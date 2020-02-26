package main

import (
	"fmt"
	"log"

	ui18n "github.com/unknwon/i18n"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// Telegram group ID consts
const (
	tAnonBalkan  = -1001161265502
	tAnon        = -1001361489843
	tAnonTaxi    = -1001422544298
	tAnonTaxiPrv = -1001271198034
	tAnonOps     = -297434742
)

func initBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(conf.TelegramAPIKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(tAnonOps, "AnonFounder successfully started. 🚀")
	bot.Send(msg)

	return bot
}

func logTelegram(message string) {
	msg := tgbotapi.NewMessage(tAnonOps, message)
	bot.Send(msg)
}

func messageTelegram(message string, groupID int64) {
	msg := tgbotapi.NewMessage(groupID, message)
	bot.Send(msg)
}

func sendGroupsMessageInvestment(investment float64) {
	msg := tgbotapi.NewMessage(tAnonTaxi, fmt.Sprintf(ui18n.Tr(lang, "newPurchase"), investment))
	bot.Send(msg)
}

func sendGroupsMessagePrice(newPrice float64) {
	msg := tgbotapi.NewMessage(tAnonTaxi, fmt.Sprintf(ui18n.Tr(lang, "priceRise"), newPrice))
	bot.Send(msg)
}

// TelegramUpdate struct represent webhook update data from Telegram
type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID                          int    `json:"id"`
			Title                       string `json:"title"`
			Type                        string `json:"type"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
		} `json:"chat"`
		Date           int `json:"date"`
		ReplyToMessage struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID        int    `json:"id"`
				IsBot     bool   `json:"is_bot"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
			} `json:"from"`
			Chat struct {
				ID                          int    `json:"id"`
				Title                       string `json:"title"`
				Type                        string `json:"type"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"reply_to_message"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
		NewChatParticipant struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"new_chat_participant"`
		NewChatMember struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"new_chat_member"`
		NewChatMembers []struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"new_chat_members"`
	} `json:"message"`
}
