package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

type Tweet struct {
	Id         int
	User_id    int
	Message    sql.NullString
	Created_at time.Time
}

func tweet_create(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	tweet_text := req.PostFormValue("tweet")

	// Get the "session_uid" cookie to get the current loged in users id
	user_id, err := req.Cookie("session_uid")

	currentTime := time.Now()

	_, err = Db.Exec("insert into tweets (user_id, msg, created_at, updated_at) values ($1, $2, $3, $3)", user_id.Value, tweet_text, currentTime)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 302)
}

func tweet_delete(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	session := Session{}
	session_cookie, err := req.Cookie("session")

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select id, user_id, token, created_at, updated_at from sessions where token = $1", session_cookie.Value).Scan(&session.Id, &session.User_id, &session.Token, &session.Created_at, &session.Updated_at)

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
