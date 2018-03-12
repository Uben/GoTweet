package main

import (
	// "database/sql"
	// "gowebapp/models"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

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
}

func delete_user_follow(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	follower_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	url_params := mux.Vars(req)
	following_id := url_params["user_id"]

	_, err = Db.Exec("delete from user_follows where follower_id = $1 AND following_id = $2", follower_id.Value, following_id)

	if err != nil {
		panic(err)
	}
}
