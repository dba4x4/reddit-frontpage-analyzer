package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/src/util"
)

var tag = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Tag",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"confidence": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var post = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Post",
		Fields: graphql.Fields{
			"author": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"domain": &graphql.Field{
				Type: graphql.String,
			},
			"subreddit": &graphql.Field{
				Type: graphql.String,
			},
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"permalink": &graphql.Field{
				Type: graphql.String,
			},
			"selftext": &graphql.Field{
				Type: graphql.String,
			},
			"thumbnail": &graphql.Field{
				Type: graphql.String,
			},
			"created_utc": &graphql.Field{
				Type: graphql.Float,
			},
			"num_comments": &graphql.Field{
				Type: graphql.Int,
			},
			"score": &graphql.Field{
				Type: graphql.Int,
			},
			"ups": &graphql.Field{
				Type: graphql.Int,
			},
			"downs": &graphql.Field{
				Type: graphql.Int,
			},
			"over_18": &graphql.Field{
				Type: graphql.Int,
			},
			"is_self": &graphql.Field{
				Type: graphql.Int,
			},
			"post_hint": &graphql.Field{
				Type: graphql.String,
			},
			"tags": &graphql.Field{
				Type: graphql.NewList(tag),
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type:        graphql.NewList(post),
				Description: "List of posts",
				Args: graphql.FieldConfigArgument{
					"date": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					date := p.Args["date"]
					db := util.InitDatabase()
					posts := []util.Post{}
					err := db.
						Preload("Tags").
						Where("date(to_timestamp(date_created)) = ? AND post_hint = 'image'", date).
						Find(&posts).
						Error
					if err != nil {
						return nil, errors.New("Internal server error")
					}
					return posts, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

type graphqlQuery struct {
	Query string `json:"query"`
}

func queryGet(w http.ResponseWriter, r *http.Request) {
	query(w, r.URL.Query().Get("query"))
}

func queryPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var graphqlQuery graphqlQuery
	err := decoder.Decode(&graphqlQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}
	query(w, graphqlQuery.Query)
}

func query(w http.ResponseWriter, query string) {
	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}
	result := executeQuery(query, schema)
	json.NewEncoder(w).Encode(result)
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
