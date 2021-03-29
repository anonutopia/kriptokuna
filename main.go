package main

import (
	"log"

	macaron "gopkg.in/macaron.v1"
	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var db *gorm.DB

var pc *PriceClient

var wm *WavesMonitor

var bot *telebot.Bot

var m *macaron.Macaron

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	bot = initTelegramBot()

	db = initDb()

	pc = initPriceClient()

	m = initMacaron()

	m.Get("/:address/interest.js", accumulatedInterest)

	logTelegram("Kriptokuna daemon successfully started. 🚀")

	initWavesMonitor()
}
