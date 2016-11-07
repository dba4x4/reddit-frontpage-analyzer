package database

import (
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
