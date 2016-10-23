package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/swordbeta/reddit-frontpage-analyzer-go/util"
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
				Type: graphql.Int,
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
			"post": &graphql.Field{
				Type:        graphql.NewList(post),
				Description: "List of posts",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					db := util.InitDatabase()
					posts := []util.Post{}
					dbf := db.Preload("Tags").Find(&posts)
					if dbf.Error != nil {
						return nil, nil
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

func query(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("query") == "" {
		return
	}
	result := executeQuery(r.URL.Query()["query"][0], schema)
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
