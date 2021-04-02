package main

import (
	"log"

	"gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var conf *Config

var db *gorm.DB

var pc *PriceClient

var wm *WavesMonitor

var bot *telebot.Bot

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	bot = initTelegramBot()

	db = initDb()

	pc = initPriceClient()

	logTelegram("Kriptokuna daemon successfully started. 🚀")

	initWavesMonitor()
}
