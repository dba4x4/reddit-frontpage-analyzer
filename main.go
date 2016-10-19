package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	if viper.GetBool("debug") {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	run()
	ch := gocron.Start()
	gocron.Every(30).Minutes().Do(run)
	<-ch
}

func run() {
	log.Println("Starting to process all posts on r/all...")
	db := initDatabase()
	defer db.Close()
	posts, err := getPosts("all")
	if err != nil {
		if err.Error() == "Too many requests" {
			log.Println("Hit reddit late limitting, waiting till next batch...")
			return
		}
		log.Fatalln(err)
	}
	var wg sync.WaitGroup
	for _, post := range posts {
		wg.Add(1)
		go processPost(post, db, &wg)
	}
	wg.Wait()
	log.Println("Finished processing all posts, waiting 30 minutes...")
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/reddit-frontpage-analyzer/")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

func initDatabase() *gorm.DB {
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

func processPost(post *Post, db *gorm.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	if alreadySaved(post.ID, db) {
		log.Println("Skipping #", post.ID, "...")
		return
	}
	log.Println("Started processing #", post.ID, "...")
	if post.PostHint == "image" {
		post.Tags = tagImg(post.URL)
	}
	savePost(post, db)
	log.Println("Finished processing #", post.ID, "...")
}
