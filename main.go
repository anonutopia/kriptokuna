package main

import (
	"github.com/anonutopia/gowaves"
	"github.com/go-macaron/binding"
	"github.com/jinzhu/gorm"
	macaron "gopkg.in/macaron.v1"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var conf *Config

var wnc *gowaves.WavesNodeClient

var db *gorm.DB

var bot *tgbotapi.BotAPI

var m *macaron.Macaron

var pc *PriceClient

var token *Token

func main() {
	conf = initConfig()

	db = initDb()

	wnc = initWaves()

	bot = initBot()

	pc = initPriceClient()

	token = initToken()

	m = initMacaron()
	m.Post("/", binding.Json(TelegramUpdate{}), pageView)

	initMonitor()
}
