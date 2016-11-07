package database

import (
	"github.com/jinzhu/gorm"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/domain"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

// TearDown tears down test tables.
func TearDown() {
	util.InitConfig()
	db := InitDatabase()
	defer db.Close()
	db.Delete(domain.Post{})
	db.Delete(domain.Tag{})
}

// LoadTestPosts saves posts to the test database.
func LoadTestPosts(db *gorm.DB) {
	SavePost(&domain.Post{
		ID:          "firstPost",
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
		ID:          "secondPost",
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
}
