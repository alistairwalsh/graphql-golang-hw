# GraphQL server in Go - hello world
If you are new to GraphQL in Go, this repository should get you started quickly. The script `Server.go` starts a http server that accepts GraphQL queries at the endpoint `/graphql`. 

I use [this](https://github.com/graphql-go/graphql) implementation of GraphQL, just because there is more documentation available. 

## Getting started
I dockerized the server to make it as easy as possible to play with the code.

    go get github.com/stefanhengl/graphql-goalang-hw
    docker build -t graphql-server .
    docker run -p 9090:9090 graphql-server

    🌍 Listening on http://localhost:9090

Of course, you can also run the server outside docker. Make sure you have all dependencies installed though.

    go run server.go
    
    🌍 Listening on http://localhost:9090

## Sending requests
    
    GET /graphql?query={hello}

    {
        "data": {
            "hello": "world"
        }
    }

    GET /graphql?query={user{first,last}}

    {
        "data": {
            "user": {
                "first": "john",
                "last": "doe"
            }
        }
    }

    GET /graphql?query=mutation {changeUser(first: "foo", last: "bar")}

    {
        "data": {
            "changeUser": true
        }
    }

    GET /graphql?query={user{first,last}}

    {
        "data": {
            "user": {
                "first": "foo",
                "last": "bar"
            }
        }
    }

## Postman
The Postman collection `graphql-golang.hw.postlam_collection` contains templates for all the requests above ([doc](https://www.getpostman.com/docs/postman/collections/data_formats)). 