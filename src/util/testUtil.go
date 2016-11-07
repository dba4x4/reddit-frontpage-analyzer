package util

import "github.com/jinzhu/gorm"

// TearDown tears down test tables.
func TearDown() {
	InitConfig()
	db := InitDatabase()
	defer db.Close()
	db.Delete(Post{})
	db.Delete(Tag{})
}

// LoadTestPosts saves posts to the test database.
func LoadTestPosts(db *gorm.DB) {
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
}
