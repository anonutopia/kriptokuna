package main

import (
	"github.com/jinzhu/gorm"
)

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `sql:"size:255;unique_index"`
	ValueInt uint64 `sql:"type:int"`
	ValueStr string `sql:"type:string"`
}

// Transaction represents node's transaction
type Transaction struct {
	gorm.Model
	TxID      string `sql:"size:255"`
	Processed bool   `sql:"DEFAULT:false"`
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address          string `sql:"size:255;unique_index"`
	TelegramUsername string `sql:"size:255"`
	TelegramID       int    `sql:"unique_index"`
	ReferralID       uint
	Referral         *User
}
