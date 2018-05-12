package main

import (
	"github.com/jinzhu/gorm"
)

const (
	DBNAME = "kkn.db"

	TYPE_BLOCKED    = 1
	TYPE_DRIVER     = 2
	TYPE_USER       = 3
	TYPE_INTERESTED = 4
)

type Hacktivist struct {
	gorm.Model
	Email string `sql:"size:255"`
	Type  int
}

type Hacker struct {
	gorm.Model
	Email string `sql:"size:255"`
	Type  string `sql:"size:255"`
}
