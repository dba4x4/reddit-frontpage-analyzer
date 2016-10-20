package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	vision "github.com/ahmdrz/microsoft-vision-golang"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Test_alreadySaved(t *testing.T) {
	initConfig()
	db := initDatabase()
	defer db.Close()
	savePost(&post{
		ID: "myExistingPost",
	}, db)
	type args struct {
		id string
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"New post",
			args{
				"myNewPost",
				db,
			},
			false,
		},
		{
			"Existing post",
			args{
				"myExistingPost",
				db,
			},
			true,
		},
	}
	for _, tt := range tests {
		if got := alreadySaved(tt.args.id, tt.args.db); got != tt.want {
			t.Errorf("%q. alreadySaved() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_savePost(t *testing.T) {
	initConfig()
	db := initDatabase()
	defer db.Close()
	post := &post{
		ID: "myPost",
	}
	if alreadySaved(post.ID, db) {
		t.Errorf("Test_savePost ID %v already exists.", post.ID)
	}
	savePost(post, db)
	if alreadySaved(post.ID, db) == false {
		t.Errorf("Test_savePost ID %v was not saved", post.ID)
	}
}

const firstImage = "http://test.com/firstimage.jpg"
const secondImage = "http://test.com/secondImage.jpg"

type MockedVision struct {
}

func (mV MockedVision) Tag(url string) (vision.VisionResult, error) {
	switch url {
	case firstImage:
		return vision.VisionResult{
			Tags: []vision.Tag{
				vision.Tag{
					Name:       "Person",
					Confidence: 0.95,
				},
			},
		}, nil
	case secondImage:
		return vision.VisionResult{
			Tags: []vision.Tag{
				vision.Tag{
					Name:       "Dog",
					Confidence: 0.95,
				},
				vision.Tag{
					Name:       "Grass",
					Confidence: 0.75,
				},
			},
		}, nil
	}
	return vision.VisionResult{}, errors.New("Could not fetch image!")
}

func Test_tagImg(t *testing.T) {
	type args struct {
		url    string
		vision tagger
	}
	tests := []struct {
		name string
		args args
		want []tag
	}{
		{
			"Person image",
			args{
				firstImage,
				MockedVision{},
			},
			[]tag{
				tag{
					Tag: vision.Tag{
						Name:       "Person",
						Confidence: 0.95,
					},
				},
			},
		},
		{
			"Dog running on grass image",
			args{
				secondImage,
				MockedVision{},
			},
			[]tag{
				tag{
					Tag: vision.Tag{
						Name:       "Dog",
						Confidence: 0.95,
					},
				},
				tag{
					Tag: vision.Tag{
						Name:       "Grass",
						Confidence: 0.75,
					},
				},
			},
		},
		{
			"Unknown image",
			args{
				"http://test.com/unknown.jpg",
				MockedVision{},
			},
			[]tag{},
		},
	}
	for _, tt := range tests {
		if got := tagImg(tt.args.url, tt.args.vision); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. tagImg() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_getPosts(t *testing.T) {
	mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, jsonTestData)
	}))
	defer mockedServer.Close()
	redditURL = mockedServer.URL + "/r/%s.json"
	res, err := getPosts("all")
	if err != nil {
		t.Error("Could not fetch reddit posts!")
	}
	if len(res) != 25 {
		t.Errorf("Reddit response did not contain 25 posts but contained %v posts", len(res))
	}
	if res[10].PostHint != "image" {
		t.Errorf("Tenth post should be an image!")
	}
}

func Test_getPostsTooManyRequests(t *testing.T) {
	mockedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer mockedServer.Close()
	redditURL = mockedServer.URL + "/r/%s.json"
	res, err := getPosts("all")
	if res != nil {
		t.Error("The webserver returned posts when it should responde with a status 'Too many requests'.")
	}
	if err != nil && err.Error() != "Too many requests" {
		t.Errorf("The webserver errored with %s, but should have errored with 'Too many requests'.", err.Error())
	}
}
