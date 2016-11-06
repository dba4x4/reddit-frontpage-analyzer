package resource

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

// GetPosts returns posts filtered by date.
func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("date") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request - Missing date querystring."))
		return
	}
	db := util.InitDatabase()
	defer db.Close()
	posts, err := util.GetPostsByDate(r.URL.Query().Get("date"), db)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "json")
	json.NewEncoder(w).Encode(posts)
}
