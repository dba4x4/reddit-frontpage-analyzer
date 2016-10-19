package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	vision "github.com/ahmdrz/microsoft-vision-golang"
	"github.com/jinzhu/gorm"
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

func getPosts(subreddit string) ([]*post, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://www.reddit.com/r/%s.json", subreddit)
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
