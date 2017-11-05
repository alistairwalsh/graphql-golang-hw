// Hello world example for a GraphQL server in golang
//
// The code starts an http server which exposes a GraphQL API
// supporting queries and mutations.
//
// Postman GET Requests:
// Query:
// http://localhost:9090/graphql?query={hello}
// http://localhost:9090/graphql?query={user{first, last}}
//
// Mutation:
// http://localhost:9090/graphql?query=mutation {changeUser(first: "foo", last: "bar")}
//
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

type user struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// Usually, the user data would be stored in a database. I define a
// global variable here to keep the code as self-contained as possible
var u = user{"john", "doe"}

func init() {
	// init initializes the  graphql schema with supported queries and mutations
	// The schema is required for parsing the query strings.
	// The query supports two fields: A simple string field "hello" and a more complex field "user"
	// The mutation supports updating the user field
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

	// Defined supported query fields. Each field comes with its own resovler
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
				return u, nil
			},
		},
	}

	// Define supported mutations
	mutations := graphql.Fields{
		"changeUser": &graphql.Field{
			Type: graphql.Boolean,
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"last": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				first := p.Args["first"].(string)
				last := p.Args["last"].(string)

				u.First = first
				u.Last = last

				return true, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation)}
	schema, _ = graphql.NewSchema(schemaConfig)
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	// read out query string
	// https://golang.org/pkg/net/http/#Request
	// https://golang.org/pkg/net/url/#URL.Query
	url := r.URL.Query()
	fmt.Println(url["query"][0])
	params := graphql.Params{Schema: schema, RequestString: url["query"][0]}

	// send query string to GraphQL
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

	// this line connects to the server to GraphQL
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
