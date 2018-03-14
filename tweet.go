package main

import (
	// "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gowebapp/models"
	"net/http"
	"time"
)

func tweet_create(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	tweet_text := req.PostFormValue("tweet")
	retUser := Models.User{}
	// Get the "session_uid" cookie to get the current loged in users id
	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

	if err != nil {
		panic(err)
	}

	current_time := time.Now()
	is_retweet := false

	_, err = Db.Exec("insert into tweets (user_id, msg, name, username, is_retweet, origin_user_id, origin_name, origin_username, created_at) values ($1, $2, $3, $4, $5, $1, $3, $4, $6)", retUser.Id, tweet_text, retUser.Name, retUser.Username, is_retweet, current_time)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func tweet_delete(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	session := Models.Session{}
	session_cookie, err := req.Cookie("session")

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select id, user_id, token, created_at from sessions where token = $1", session_cookie.Value).Scan(&session.Id, &session.User_id, &session.Token, &session.Created_at)

	if err == nil {
		url_params := mux.Vars(req)
		tweet_id := url_params["tweet_id"]
		tweet := Models.Tweet{}

		err = Db.QueryRow("select id, user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at from tweets where id = $1", tweet_id).Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

		if err != nil {
			panic(err)
		}

		if tweet.User_id == session.User_id {
			fmt.Printf("\nDeleting the tweet id of %s...", tweet_id)
			_, err := Db.Exec("delete from tweets where id = $1", tweet_id)
			fmt.Printf("Done.\n")

			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func retweet_create(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	tweet := Models.Tweet{}
	retUser := Models.User{}
	is_retweet := true
	current_time := time.Now()

	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]

	user_id, err := req.Cookie("session_uid")

	err = Db.QueryRow("select id, user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at from tweets where id = $1 limit 1", tweet_id).Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select id, name, username from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Name, &retUser.Username)

	if err != nil {
		panic(err)
	}

	if tweet.Is_retweet == false {
		_, err = Db.Exec("insert into tweets (user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", retUser.Id, tweet.Message, retUser.Name, retUser.Username, is_retweet, tweet.Id, tweet.User_id, tweet.Name, tweet.Username, current_time)
	} else {
		_, err = Db.Exec("insert into tweets (user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", retUser.Id, tweet.Message, retUser.Name, retUser.Username, is_retweet, tweet.Otweet_id, tweet.Ouser_id, tweet.Oname, tweet.Ousername, current_time)
	}

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func favorite_tweet(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]
	current_time := time.Now()

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("insert into favorites (user_id, tweet_id, created_at) values ($1, $2, $3)", user_id.Value, tweet_id, current_time)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func unfavorite_tweet(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	_, err = Db.Exec("delete from favorites where user_id = $1 and tweet_id = $2", user_id.Value, tweet_id)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

// func is_favorite(tweet_id, user_id int) bool {

// 	fave := Favorite{}
// 	err := Db.QueryRow("select id, user_id, tweet_id, created_at from favorites where user_id = $1 and tweet_id = $2", strconv.Itoa(user_id), strconv.Itoa(tweet_id)).Scan(&fave.Id, &fave.User_id, &fave.Tweet_id, &fave.Created_at)

// 	if err == sql.ErrNoRows {
// 		return false
// 	} else if err != nil {
// 		panic(err)
// 	}

// 	return true
// }
