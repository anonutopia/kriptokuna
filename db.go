package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func initDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "kriptokuna.db")

	if err != nil {
		log.Printf("[initDb] error: %s", err)
	}

	db.DB()
	db.DB().Ping()
	db.LogMode(conf.Debug)
	db.AutoMigrate(&KeyValue{}, &Transaction{}, &User{})

	return db
}
