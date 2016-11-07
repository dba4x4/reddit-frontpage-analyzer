package analyzer

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	vision "github.com/ahmdrz/microsoft-vision-golang"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

const firstImage = "http://test.com/firstimage.jpg"
const secondImage = "http://test.com/secondImage.jpg"

type mockedVision struct {
}

func (mV mockedVision) Tag(url string) (vision.VisionResult, error) {
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

func TestMain(m *testing.M) {
	ret := m.Run()
	util.TearDown()
	os.Exit(ret)
}

func Test_tagImg(t *testing.T) {
	type args struct {
		url    string
		vision util.Tagger
	}
	tests := []struct {
		name string
		args args
		want []util.Tag
	}{
		{
			"Person image",
			args{
				firstImage,
				mockedVision{},
			},
			[]util.Tag{
				util.Tag{
					Name:       "Person",
					Confidence: 0.95,
				},
			},
		},
		{
			"Dog running on grass image",
			args{
				secondImage,
				mockedVision{},
			},
			[]util.Tag{
				util.Tag{
					Name:       "Dog",
					Confidence: 0.95,
				},
				util.Tag{
					Name:       "Grass",
					Confidence: 0.75,
				},
			},
		},
		{
			"Unknown image",
			args{
				"http://test.com/unknown.jpg",
				mockedVision{},
			},
			[]util.Tag{},
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

func Test_processPost(t *testing.T) {
	post := &util.Post{
		ID: "processPost",
	}
	util.InitConfig()
	db := util.InitDatabase()
	defer db.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	if processPost(post, db, &mockedVision{}, &wg) != true {
		t.Error("Failed to process a new post, it already is processed!")
	}
}

func Test_processPostAlreadyProcessed(t *testing.T) {
	post := &util.Post{
		ID: "existingPost",
	}
	util.InitConfig()
	db := util.InitDatabase()
	defer db.Close()
	util.SavePost(post, db)
	var wg sync.WaitGroup
	wg.Add(1)
	if processPost(post, db, &mockedVision{}, &wg) != false {
		t.Error("Failed to process a new post, it already is processed!")
	}
}
