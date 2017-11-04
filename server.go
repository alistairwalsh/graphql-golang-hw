// http://mycodesmells.com/post/building-graphql-api-in-go
// http://mycodesmells.com/post/advanced-graphql-in-golang
// https://github.com/mycodesmells/graphql-example
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
)

func graphqlHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Query()

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
	schema, _ := graphql.NewSchema(schemaConfig)

	params := graphql.Params{Schema: schema, RequestString: url["query"][0]}
	res := graphql.Do(params)
	if res.HasErrors() {
		log.Fatalf("Failed due to errors: %v\n", res.Errors)
	}

	rJSON, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(rJSON)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/graphql", graphqlHandler)
	http.Handle("/", r)

	// set port and inform the user
	port := 9090
	fmt.Printf("🌍 Listening on http://localhost:%d", port)

	// start server
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
