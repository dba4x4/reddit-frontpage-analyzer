package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type response struct {
	Data struct {
		Children []struct {
			Data *Post
		}
	}
}

func getPosts(subreddit string) []*Post {
	client := &http.Client{}
	url := fmt.Sprintf("https://www.reddit.com/r/%s.json", subreddit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "reddit-frontpage-analyzer-go:v0.1.0 (by /u/swordbeta)")
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusTooManyRequests {
		log.Fatalln("Too many requests!")
	}
	r := new(response)
	err = json.NewDecoder(res.Body).Decode(r)
	if err != nil {
		log.Fatalln(err)
	}
	posts := make([]*Post, len(r.Data.Children))
	for i, child := range r.Data.Children {
		posts[i] = child.Data
	}
	return posts
}
