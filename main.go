package main

import (
	"log"

	"gorm.io/gorm"
)

var conf *Config

var db *gorm.DB

var pc *PriceClient

var wm *WavesMonitor

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conf = initConfig()

	db = initDb()

	pc = initPriceClient()

	initWavesMonitor()
}
