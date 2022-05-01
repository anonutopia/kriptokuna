package main

import "gorm.io/gorm"

// KeyValue model is used for storing key/values
type KeyValue struct {
	gorm.Model
	Key      string `gorm:"size:255;uniqueIndex"`
	ValueInt uint64 `gorm:"type:int"`
	ValueStr string `gorm:"type:string"`
}

// Transaction represents node's transaction
type Transaction struct {
	gorm.Model
	TxID      string `gorm:"size:255;uniqueIndex"`
	Processed bool   `gorm:"DEFAULT:false"`
}

// User represents Telegram user
type User struct {
	gorm.Model
	Address        string `gorm:"size:255;uniqueIndex"`
	AmountAhrk     uint
	AmountAhrkAint uint
	AmountAeur     uint
	AmountAeurAint uint
	AmountWaves    uint
	ReferralID     uint
	Referral       *User
}
