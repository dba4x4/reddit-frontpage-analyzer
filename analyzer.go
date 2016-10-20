package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	vision "github.com/ahmdrz/microsoft-vision-golang"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

type post struct {
	Author       string  `json:"author"`
	Title        string  `json:"title"`
	URL          string  `json:"url"`
	Domain       string  `json:"domain"`
	Subreddit    string  `json:"subreddit"`
	ID           string  `json:"id" gorm:"primary_key"`
	Permalink    string  `json:"permalink"`
	Selftext     string  `json:"selftext"`
	ThumbnailURL string  `json:"thumbnail"`
	DateCreated  float32 `json:"created_utc"`
	NumComments  int     `json:"num_comments"`
	Score        int     `json:"score"`
	Ups          int     `json:"ups"`
	Downs        int     `json:"downs"`
	IsNSFW       bool    `json:"over_18"`
	IsSelf       bool    `json:"is_self"`
	PostHint     string  `json:"post_hint"`
	Tags         []tag   `json:"-"`
}

type tag struct {
	ID     uint `gorm:"primary_key"`
	PostID string
	vision.Tag
}

type tagger interface {
	Tag(url string) (vision.VisionResult, error)
}

type response struct {
	Data struct {
		Children []struct {
			Data *post
		}
	}
}

var redditURL = "https://www.reddit.com/r/%s.json"

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
	vision, err := vision.New(viper.GetString("microsoft.key"))
	if err != nil {
		log.Fatalln(err)
	}
	var wg sync.WaitGroup
	for _, post := range posts {
		wg.Add(1)
		go processPost(post, db, vision, &wg)
	}
	wg.Wait()
	log.Println("Finished processing all posts, waiting 30 minutes...")
}

func initConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
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
	db.AutoMigrate(&post{}, &tag{})
	return db
}

func processPost(post *post, db *gorm.DB, vision *vision.Vision, wg *sync.WaitGroup) {
	defer wg.Done()
	if alreadySaved(post.ID, db) {
		log.Println("Skipping #", post.ID, "...")
		return
	}
	log.Println("Started processing #", post.ID, "...")
	if post.PostHint == "image" {
		post.Tags = tagImg(post.URL, vision)
	}
	savePost(post, db)
	log.Println("Finished processing #", post.ID, "...")
}

func getPosts(subreddit string) ([]*post, error) {
	client := &http.Client{}
	url := fmt.Sprintf(redditURL, subreddit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "reddit-frontpage-analyzer-go:v0.1.0 (by /u/swordbeta)")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("Too many requests")
	}
	r := new(response)
	err = json.NewDecoder(res.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	posts := make([]*post, len(r.Data.Children))
	for i, child := range r.Data.Children {
		posts[i] = child.Data
	}
	return posts, nil
}

func alreadySaved(id string, db *gorm.DB) bool {
	dbf := db.Where("id = ?", id).First(&post{})
	return dbf.Error == nil
}

func savePost(post *post, db *gorm.DB) {
	dbc := db.Create(&post)
	if dbc.Error != nil {
		log.Fatalln(dbc.Error)
	}
}

func tagImg(url string, vision tagger) []tag {
	result, err := vision.Tag(url)
	if err != nil {
		log.Println(fmt.Sprintf("While trying to tag %s got the following error: %s", url, err))
		return make([]tag, 0)
	}
	response := make([]tag, len(result.Tags))
	for i, visionTag := range result.Tags {
		response[i] = tag{
			Tag: visionTag,
		}
	}
	return response
}
