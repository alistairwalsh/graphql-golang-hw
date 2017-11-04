// Hello world example for a graphql API exposed through an http server
// For more information see the following links:
//
// Format of GET request
// http://graphql.org/learn/serving-over-http/#get-request
//
// Big parts of the code have been inspired by
// http://mycodesmells.com/post/building-graphql-api-in-go
// http://mycodesmells.com/post/advanced-graphql-in-golang
// https://github.com/mycodesmells/graphql-example
//
// TODO: http://mycodesmells.com/post/making-changes-in-graphql-api
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

var schema graphql.Schema

func init() {
	// init initializes the  graphql schema with two fields:
	// A simple string field "hello" and a more complex field "user"
	var userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"first": &graphql.Field{
				Type: graphql.String,
			},
			"last": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	type User struct {
		First string `json:"first"`
		Last  string `json:"last"`
	}

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"user": &graphql.Field{
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var u = User{"john", "doe"}
				return u, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ = graphql.NewSchema(schemaConfig)
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {

	// read out query string
	// https://golang.org/pkg/net/http/#Request
	// https://golang.org/pkg/net/url/#URL.Query
	url := r.URL.Query()
	params := graphql.Params{Schema: schema, RequestString: url["query"][0]}

	// send query string to graphql
	res := graphql.Do(params)
	if res.HasErrors() {
		log.Fatalf("Failed due to errors: %v\n", res.Errors)
	}

	// convert response to json
	rJSON, _ := json.Marshal(res)

	// return to client
	// http://www.alexedwards.net/blog/golang-response-snippets#json
	w.Header().Set("Content-Type", "application/json")
	w.Write(rJSON)
}

func main() {
	// create route
	r := mux.NewRouter()
	r.HandleFunc("/graphql", graphqlHandler)
	http.Handle("/", r)

	// set port and inform the user
	port := 9090
	fmt.Printf("🌍 Listening on http://localhost:%d\n", port)

	// start server
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
