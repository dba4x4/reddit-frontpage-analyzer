package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

func alreadySaved(id string, db *gorm.DB) bool {
	dbf := db.Where("id = ?", id).First(&Post{})
	return dbf.Error == nil
}

func savePost(post *Post, db *gorm.DB) {
	dbc := db.Create(&post)
	if dbc.Error != nil {
		log.Fatalln(dbc.Error)
	}
}
