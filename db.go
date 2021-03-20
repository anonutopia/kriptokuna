package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDb() *gorm.DB {
	dbconf := gorm.Config{}

	if conf.Debug {
		dbconf.Logger = logger.Default.LogMode(logger.Info)
	} else {
		dbconf.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(sqlite.Open("kriptokuna.db"), &dbconf)

	if err != nil {
		log.Println(err)
	}

	if err := db.AutoMigrate(&KeyValue{}, &Transaction{}); err != nil {
		panic(err.Error())
	}

	return db
}
