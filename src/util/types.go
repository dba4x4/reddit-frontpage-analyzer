package util

import vision "github.com/ahmdrz/microsoft-vision-golang"

// Post represents a post on reddit
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
	Tags         []Tag   `json:"tags"`
}

// Tag represents a Microsoft Vision API.
type Tag struct {
	ID         uint `gorm:"primary_key"`
	PostID     string
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

// Tagger is an API that tags an image.
type Tagger interface {
	Tag(url string) (vision.VisionResult, error)
}
