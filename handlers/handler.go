package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Msg string
}

type GoogleHandler struct {
	RedisConfig *redis.Options
	Oauth *oauth2.Config
	Db *sql.DB
}

func (gHandler GoogleHandler) InsertIntoMySQL(id,email,name,access_token string) error{
	db := gHandler.Db
	if db == nil {
		return errors.New("DB is nil")
	}
	insertQuery, err := db.Prepare("INSERT INTO user(id,email,name,access_token) values(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = insertQuery.Exec(id,email,name,access_token)

	if err != nil {
		return err
	}
	return nil
}

func (gHandler GoogleHandler) OauthLogin(res http.ResponseWriter,req *http.Request){
	url := gHandler.Oauth.AuthCodeURL("hello")
	http.Redirect(res,req,url,http.StatusTemporaryRedirect)
}

func (gHandler GoogleHandler) OauthCallback(res http.ResponseWriter,req *http.Request){

	res.Header().Add("Content-Type","application/json")
	state := req.FormValue("state")
	if state != "hello"{
		result := Response{Msg: "State did not match"}
		json.NewEncoder(res).Encode(&result)
		return
	}

	token,err := gHandler.Oauth.Exchange(oauth2.NoContext,req.FormValue("code"))
	if err != nil {
		res.Write([]byte("Error occured 2"+err.Error()))
		return
	}


	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		result := Response{Msg: err.Error()}
		json.NewEncoder(res).Encode(&result)
		return
	}
	var v interface{}
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body,&v)
	userMap:= v.(map[string]interface{})
	ctx := context.Background()
	access_token := token.AccessToken
	name := userMap["name"]
	email := userMap["email"]
	id := userMap["id"]
	client := redis.NewClient(gHandler.RedisConfig)
	_, err = client.Get(ctx,email.(string)).Result()

	if err != nil {
		log.Println("Stored to Database and cache")
		err = gHandler.InsertIntoMySQL(id.(string),email.(string),name.(string),access_token)
		_,err = client.Set(ctx,email.(string),access_token,0).Result()
	} else {
		log.Println("Served from cache")
	}

	json.NewEncoder(res).Encode(v)
}