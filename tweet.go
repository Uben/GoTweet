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

	// Get the "session_uid" cookie to get the current loged in users id
	user_id, err := req.Cookie("session_uid")

	current_time := time.Now()

	_, err = Db.Exec("insert into tweets (user_id, msg, created_at) values ($1, $2, $3)", user_id.Value, tweet_text, current_time)

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

		fmt.Printf("\nDeleting the tweet id of %s...", tweet_id)
		_, err := Db.Exec("delete from tweets where id = $1", tweet_id)
		fmt.Printf("Done.\n")

		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func retweet_create(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func retweet_delete(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

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
