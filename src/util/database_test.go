package util

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	TearDown()
	os.Exit(ret)
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
