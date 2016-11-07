package resource

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util/database"
)

// GetPosts returns posts filtered by date.
func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("date") == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request - Missing date querystring."))
		return
	}
	db := database.InitDatabase()
	defer db.Close()
	posts, err := database.GetPostsByDate(r.URL.Query().Get("date"), db)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
