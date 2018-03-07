package main

import (
	// "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

type Follow struct {
	Id           int
	follower_id  int
	following_id int
	Created_at   time.Time
	Updated_at   time.Time
}

func create_user_follow(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	follower_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	url_params := mux.Vars(req)
	following_id := url_params["user_id"]

	current_time := time.Now()

	_, err = Db.Exec("insert into user_follows (follower_id, following_id, created_at, updated_at) values ($1, $2, $3, $3)", follower_id.Value, following_id, current_time)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 200)
}

func delete_user_follow(res http.ResponseWriter, req *http.Request) {

}
