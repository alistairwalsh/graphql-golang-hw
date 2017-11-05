# golang graphql server - hello world

## Getting started
    go get github.com/stefanhengl/graphql-goalang-hw
    docker build -t graphql-server .
    docker run graphql-server

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
Open the Postman collection `graphql-golang.hw.postlam_collection` which contains templates for all the requests above.