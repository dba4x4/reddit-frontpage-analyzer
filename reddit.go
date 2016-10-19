package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type response struct {
	Data struct {
		Children []struct {
			Data *Post
		}
	}
}

func getPosts(subreddit string) ([]*Post, error) {
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
	posts := make([]*Post, len(r.Data.Children))
	for i, child := range r.Data.Children {
		posts[i] = child.Data
	}
	return posts, nil
}
