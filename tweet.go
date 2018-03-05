package main

import (
	// "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

type Tweet struct {
	Id         int
	User_id    int
	Message    string
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
	http.Redirect(res, req, "/", 200)
}
