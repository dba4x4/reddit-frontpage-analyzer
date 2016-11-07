package database

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/domain"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	TearDown()
	os.Exit(ret)
}

func Test_PostExists(t *testing.T) {
	util.InitConfig()
	db := InitDatabase()
	defer db.Close()
	SavePost(&domain.Post{
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
		if got := PostExists(tt.args.id, tt.args.db); got != tt.want {
			t.Errorf("%q. Exists() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func Test_SavePost(t *testing.T) {
	util.InitConfig()
	db := InitDatabase()
	defer db.Close()
	post := &domain.Post{
		ID: "myPost",
	}
	if PostExists(post.ID, db) {
		t.Errorf("Test_SavePost ID %v already exists.", post.ID)
	}
	SavePost(post, db)
	if PostExists(post.ID, db) == false {
		t.Errorf("Test_SavePost ID %v was not saved", post.ID)
	}
}

func Test_GetPostsByDate(t *testing.T) {
	util.InitConfig()
	db := InitDatabase()
	defer db.Close()
	SavePost(&domain.Post{
		ID:          "firstPost1",
		DateCreated: 1451606400, // 2016-01-01
		PostHint:    "image",
		Tags: []domain.Tag{
			domain.Tag{
				Name:       "Quite sure",
				Confidence: 0.95,
			},
			domain.Tag{
				Name:       "Not quite sure",
				Confidence: 0.45,
			},
		},
	}, db)
	SavePost(&domain.Post{
		ID:          "secondPost1",
		DateCreated: 1451692800, // 2016-01-02
		PostHint:    "image",
		Tags: []domain.Tag{
			domain.Tag{
				Name:       "Quite sure",
				Confidence: 0.95,
			},
			domain.Tag{
				Name:       "Not quite sure",
				Confidence: 0.45,
			},
		},
	}, db)
	result, err := GetPostsByDate("2016-01-01", db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 1 {
		t.Errorf("Posts from 2016-01-01 had length %v instead of expected 1.", len(result))
	}
	if result[0].PostHint != "image" {
		t.Errorf("Expected PostHint to be iamge, is %v.", result[0].PostHint)
	}
	if len(result[0].Tags) != 2 {
		t.Errorf("Expected 2 tags, post has %v.", len(result[0].Tags))
	}
	result, err = GetPostsByDate("2016-01-03", db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) > 0 {
		t.Errorf("Posts from 2016-01-03 had length %v instead of expected 0.", len(result))
	}
}
