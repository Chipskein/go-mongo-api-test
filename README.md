# [go-mongo-api-test](https://go-mongo-api-test.herokuapp.com/)
<img src="https://raw.githubusercontent.com/mongodb/mongo-go-driver/master/etc/assets/mongo-gopher.png" style="width:250px"></img>
## Description
*Simple golang api to test mongodb*
## Routes
    GET /users
    GET /users/{id}
    GET /users/{id}/delete
    POST /users {
        email,
        password,
        name
    }
    POST /users/update {
        _id,
        email?,
        password?,
        name?
    }
    
## Run Locallly
#### config .env
* MongoDBURL => mongo db url
* PORT =>Server port
#### Run
        go build -o bin/go-mongo-api-test .
        ./bin/go-mongo-api-test
#### or 
        go run main.go
        
