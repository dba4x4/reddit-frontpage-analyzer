package resource

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	util.TearDown()
	os.Exit(ret)
}

func Test_GetPosts(t *testing.T) {
	util.InitConfig()
	db := util.InitDatabase()
	defer db.Close()
	util.LoadTestPosts(db)
	r, _ := http.NewRequest("GET", "/api/v1/posts", nil)
	w := httptest.NewRecorder()
	GetPosts(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Excepted 400 bad request, got %v.", w.Code)
	}
	r, _ = http.NewRequest("GET", "/api/v1/posts?date=2016-01-01", nil)
	w = httptest.NewRecorder()
	GetPosts(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v.", w.Code)
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected application/json, got %v.", contentType)
	}
	decoder := json.NewDecoder(w.Body)
	posts := []util.Post{}
	err := decoder.Decode(&posts)
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 1 || posts[0].ID != "firstPost" {
		t.Errorf(
			"Expected 1 post with ID 'firstPost', got %v posts.",
			len(posts),
		)
	} else if len(posts[0].Tags) != 2 {
		t.Errorf(
			"Expected post with ID 'firstPost' to have 2 tags, got %v.",
			len(posts[0].Tags),
		)
	}
}
