package util

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

// TestMain runs the tests and cleans up the database afterwards.
func TestMain(m *testing.M) {
	ret := m.Run()
	tearDown()
	os.Exit(ret)
}

func tearDown() {
	db := InitDatabase()
	db.Delete(Post{})
	db.Delete(Tag{})
}

func Test_Exists(t *testing.T) {
	InitConfig()
	db := InitDatabase()
	defer db.Close()
	SavePost(&Post{
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
		if got := Exists(tt.args.id, tt.args.db); got != tt.want {
			t.Errorf("%q. Exists() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_SavePost(t *testing.T) {
	InitConfig()
	db := InitDatabase()
	defer db.Close()
	post := &Post{
		ID: "myPost",
	}
	if Exists(post.ID, db) {
		t.Errorf("Test_SavePost ID %v already exists.", post.ID)
	}
	SavePost(post, db)
	if Exists(post.ID, db) == false {
		t.Errorf("Test_SavePost ID %v was not saved", post.ID)
	}
}

func Test_GetPostsByDate(t *testing.T) {
	InitConfig()
	db := InitDatabase()
	defer db.Close()
	SavePost(&Post{
		ID:          "firstPost",
		DateCreated: 1451606400, // 2016-01-01
		PostHint:    "image",
		Tags: []Tag{
			Tag{
				Name:       "Quite sure",
				Confidence: 0.95,
			},
			Tag{
				Name:       "Not quite sure",
				Confidence: 0.45,
			},
		},
	}, db)
	SavePost(&Post{
		ID:          "secondPost",
		DateCreated: 1451692800, // 2016-01-02
		PostHint:    "image",
		Tags: []Tag{
			Tag{
				Name:       "Quite sure",
				Confidence: 0.95,
			},
			Tag{
				Name:       "Not quite sure",
				Confidence: 0.45,
			},
		},
	}, db)
	result, err := GetPostsByDate("2016-01-01", db)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("Posts from 2016-01-01 had length %v instead of 1.", len(result))
	}
	if result[0].PostHint != "image" {
		t.Errorf("Post is not hinted to be an image.")
	}
	if len(result[0].Tags) != 2 {
		t.Errorf("Post should have two tags.")
	}
}
