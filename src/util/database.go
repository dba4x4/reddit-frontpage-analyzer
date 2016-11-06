package util

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	// The database connection is used in multiple packages and is
	// in this one for testing purposes.
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

// InitDatabase initializes the database.
func InitDatabase() *gorm.DB {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
			viper.GetString("postgresql.hostname"),
			viper.GetString("postgresql.username"),
			viper.GetString("postgresql.database"),
			viper.GetString("postgresql.password"),
			viper.GetString("postgresql.port"),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Post{}, &Tag{})
	return db
}

// Exists checks if a posts already exits.
func Exists(id string, db *gorm.DB) bool {
	dbf := db.Where("id = ?", id).First(&Post{})
	return dbf.Error == nil
}

// SavePost saves a new post.
func SavePost(post *Post, db *gorm.DB) {
	dbc := db.Create(&post)
	if dbc.Error != nil {
		log.Fatalln(dbc.Error)
	}
}
