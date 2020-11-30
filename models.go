package main

import (
	"github.com/jinzhu/gorm"
)

type SMS struct {
	gorm.Model
	Number string
	Text   string
	Date   int64
}
