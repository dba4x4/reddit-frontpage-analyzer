package main

import vision "github.com/ahmdrz/microsoft-vision-golang"

type Post struct {
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
	Tags         []Tag   `json:"-"`
}

type Tag struct {
	ID     uint `gorm:"primary_key"`
	PostID string
	vision.Tag
}
