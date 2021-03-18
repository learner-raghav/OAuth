package main

import (
	"Auth/driver"
	"Auth/entity"
	"Auth/handlers"
	"github.com/go-redis/redis"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback", //Need to specify this at GCP
	ClientID:     "XXXXXXXXXXXXXXXXXXXXXXX",
	ClientSecret: "XXXXXXXXXXXXXXXXXXXXXXX",
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email","https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

var redisConfig = &redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
}


var conf = entity.MySQLConfig{
	DbUser: "XXXXX",
	DbPass: "XXXXX",
	DbName: "testDB",
}


func main(){
	var db, err = driver.ConnectToDB(conf)
	if err != nil {
		log.Fatal("Could not connect to DB")
	}
	mux := http.NewServeMux()

	mux.Handle("/",http.FileServer(http.Dir("templates/")))

	oauthHandler := handlers.GoogleHandler{Oauth: googleOauthConfig,RedisConfig: redisConfig,Db: db}

	mux.HandleFunc("/auth/google/login",oauthHandler.OauthLogin)
	mux.HandleFunc("/auth/google/callback",oauthHandler.OauthCallback)
	log.Println("Server started at PORT 8000")
	if err := http.ListenAndServe(":8000",mux); err != nil {
		log.Fatal("Error occured while starting the server")
	}
}