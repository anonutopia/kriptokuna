package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/anonutopia/gowaves"
	ui18n "github.com/unknwon/i18n"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

const satInBtc = uint64(100000000)

const lang = "hr"

func executeBotCommand(tu TelegramUpdate) {
	if strings.HasPrefix(tu.Message.Text, "/price") {
		priceCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/start") {
		startCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/address") {
		addressCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/balance") {
		balanceCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/drop") {
		dropCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/freeinfo") {
		freeinfoCommand(tu)
	} else if strings.HasPrefix(tu.Message.Text, "/") {
		unknownCommand(tu)
	} else if tu.UpdateID != 0 {
		if tu.Message.ReplyToMessage.MessageID == 0 {
			if tu.Message.NewChatMember.ID != 0 &&
				!tu.Message.NewChatMember.IsBot {
				registerNewUsers(tu)
			}
		} else {
			avr, err := wnc.AddressValidate(tu.Message.Text)
			if err != nil {
				logTelegram(err.Error())
				messageTelegram(ui18n.Tr(lang, "error"), int64(tu.Message.Chat.ID))
			} else {
				if !avr.Valid {
					messageTelegram(ui18n.Tr(lang, "addressNotValid"), int64(tu.Message.Chat.ID))
				} else {
					tu.Message.Text = fmt.Sprintf("/drop %s", tu.Message.Text)
					dropCommand(tu)
				}
			}
		}
	}
}

func registerNewUsers(tu TelegramUpdate) {
	rUser := &User{TelegramID: tu.Message.From.ID}
	db.First(rUser, rUser)

	for _, user := range tu.Message.NewChatMembers {
		messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "welcome"), tu.Message.NewChatMember.FirstName), int64(tu.Message.Chat.ID))

		u := &User{TelegramID: user.ID, TelegramUsername: user.Username, ReferralID: rUser.ID}
		err := db.FirstOrCreate(u, u)
		log.Println(err)
	}
}

func priceCommand(tu TelegramUpdate) {
	kv := &KeyValue{Key: "tokenPrice"}
	db.First(kv, kv)
	price := float64(kv.ValueInt) / float64(satInBtc)
	messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "currentPrice"), price), int64(tu.Message.Chat.ID))
}

func startCommand(tu TelegramUpdate) {
	messageTelegram("Hello and welcome to Anonutopia!", int64(tu.Message.Chat.ID))
}

func addressCommand(tu TelegramUpdate) {
	messageTelegram(ui18n.Tr(lang, "myAddress"), int64(tu.Message.Chat.ID))
	messageTelegram(conf.NodeAddress, int64(tu.Message.Chat.ID))
	pc := tgbotapi.NewPhotoUpload(int64(tu.Message.Chat.ID), "qrcode.png")
	pc.Caption = "QR Code"
	bot.Send(pc)
}

func balanceCommand(tu TelegramUpdate) {
	b, err := wnc.AddressesBalance(conf.NodeAddress)
	if err != nil {
		log.Printf("balanceCommand error: %s", err)
	}
	messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "currentBalance"), float64(b.Balance)/float64(satInBtc)), int64(tu.Message.Chat.ID))
}

func dropCommand(tu TelegramUpdate) {
	kv := &KeyValue{Key: "airdropSent"}
	db.First(kv, kv)
	if kv.ValueInt >= conf.Airdrop {
		messageTelegram(ui18n.Tr(lang, "airdropFinished"), int64(tu.Message.Chat.ID))
		return
	}

	msgArr := strings.Fields(tu.Message.Text)
	if len(msgArr) == 1 && strings.HasPrefix(tu.Message.Text, "/drop") {
		msg := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), ui18n.Tr(lang, "pleaseEnter"))
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		msg.ReplyToMessageID = tu.Message.MessageID
		bot.Send(msg)
	} else {
		avr, err := wnc.AddressValidate(msgArr[1])
		if err != nil {
			logTelegram(err.Error())
			messageTelegram(ui18n.Tr(lang, "error"), int64(tu.Message.Chat.ID))
		} else {
			if !avr.Valid {
				messageTelegram(ui18n.Tr(lang, "addressNotValid"), int64(tu.Message.Chat.ID))
			} else {
				user := &User{TelegramID: tu.Message.From.ID}
				db.First(user, user)

				if len(user.Address) > 0 {
					if user.Address == msgArr[1] {
						messageTelegram(ui18n.Tr(lang, "alreadyActivated"), int64(tu.Message.Chat.ID))
					} else {
						messageTelegram(ui18n.Tr(lang, "hacker"), int64(tu.Message.Chat.ID))
					}
				} else {
					if msgArr[1] == conf.NodeAddress {
						messageTelegram(ui18n.Tr(lang, "yourAddress"), int64(tu.Message.Chat.ID))
					} else {
						atr := &gowaves.AssetsTransferRequest{
							Amount:    100000000,
							AssetID:   conf.TokenID,
							Fee:       100000,
							Recipient: msgArr[1],
							Sender:    conf.NodeAddress,
						}

						_, err := wnc.AssetsTransfer(atr)
						if err != nil {
							messageTelegram(ui18n.Tr(lang, "error"), int64(tu.Message.Chat.ID))
							logTelegram(err.Error())
						} else {
							user.TelegramID = tu.Message.From.ID
							user.TelegramUsername = tu.Message.From.Username
							user.Address = msgArr[1]
							db.Save(user)

							kv.ValueInt = kv.ValueInt + 1

							if user.ReferralID != 0 {
								rUser := &User{}
								db.First(rUser, user.ReferralID)
								if len(rUser.Address) > 0 {
									atr := &gowaves.AssetsTransferRequest{
										Amount:    50000000,
										AssetID:   conf.TokenID,
										Fee:       100000,
										Recipient: rUser.Address,
										Sender:    conf.NodeAddress,
									}

									_, err := wnc.AssetsTransfer(atr)
									if err != nil {
										logTelegram(err.Error())
									} else {
										messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "tokenSentR"), rUser.TelegramUsername), int64(tu.Message.Chat.ID))
									}
								}
							}

							db.Save(kv)

							messageTelegram(fmt.Sprintf(ui18n.Tr(lang, "tokenSent"), tu.Message.From.Username), int64(tu.Message.Chat.ID))
						}
					}
				}
			}
		}
	}
}

func freeinfoCommand(tu TelegramUpdate) {
	kv := &KeyValue{Key: "airdropSent"}
	db.First(kv, kv)
	msg := fmt.Sprintf("Tokens dropped so far: <strong>%d ATST</strong>\nTokens left to drop: <strong>%d ATST</strong>", kv.ValueInt, (conf.Airdrop - kv.ValueInt))
	log.Println(msg)
	m := tgbotapi.NewMessage(int64(tu.Message.Chat.ID), msg)
	m.ParseMode = "HTML"
	err, e := bot.Send(m)
	log.Println(err)
	log.Println(e)
}

func unknownCommand(tu TelegramUpdate) {
	messageTelegram(ui18n.Tr(lang, "commandNotAvailable"), int64(tu.Message.Chat.ID))
}
