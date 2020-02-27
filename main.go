package main

import (
	"github.com/anonutopia/gowaves"
	"github.com/jinzhu/gorm"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var conf *Config

var wnc *gowaves.WavesNodeClient

var db *gorm.DB

var bot *tgbotapi.BotAPI

// var m *macaron.Macaron

var pc *PriceClient

func main() {
	conf = initConfig()

	db = initDb()

	wnc = initWaves()

	bot = initBot()

	pc = initPriceClient()

	// m = initMacaron()
	// m.Post("/", binding.Json(TelegramUpdate{}), pageView)

	initMonitor()
}
