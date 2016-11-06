package analyzer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	vision "github.com/ahmdrz/microsoft-vision-golang"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

type response struct {
	Data struct {
		Children []struct {
			Data *util.Post
		}
	}
}

// Start starts the analyzer which runs every 30 minutes.
func Start() {
	run()
	ch := gocron.Start()
	gocron.Every(30).Minutes().Do(run)
	<-ch
}

func run() {
	log.Println("Starting to process all posts on r/all...")
	db := util.InitDatabase()
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

func processPost(post *util.Post, db *gorm.DB, vision util.Tagger, wg *sync.WaitGroup) bool {
	defer wg.Done()
	if util.Exists(post.ID, db) {
		log.Println("Skipping #", post.ID, "...")
		return false
	}
	log.Println("Started processing #", post.ID, "...")
	if post.PostHint == "image" {
		post.Tags = tagImg(post.URL, vision)
	}
	util.SavePost(post, db)
	log.Println("Finished processing #", post.ID, "...")
	return true
}

var redditURL = "https://www.reddit.com/r/%s.json"

func getPosts(subreddit string) ([]*util.Post, error) {
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
	posts := make([]*util.Post, len(r.Data.Children))
	for i, child := range r.Data.Children {
		posts[i] = child.Data
	}
	return posts, nil
}

func tagImg(url string, vision util.Tagger) []util.Tag {
	result, err := vision.Tag(url)
	if err != nil {
		log.Println(fmt.Sprintf("While trying to tag %s got the following error: %s", url, err))
		return make([]util.Tag, 0)
	}
	response := make([]util.Tag, len(result.Tags))
	for i, visionTag := range result.Tags {
		response[i] = util.Tag{
			Name:       visionTag.Name,
			Confidence: visionTag.Confidence,
		}
	}
	return response
}
