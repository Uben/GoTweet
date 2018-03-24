package controllers

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gowebapp/models"
	"log"
	"net/http"
	"time"
)

type FollowController struct {
	Db *sql.DB
}

func NewFollowController(DBCon *sql.DB) *FollowController {
	return &FollowController{DBCon}
}

func (ctrl *FollowController) Create(res http.ResponseWriter, req *http.Request) {
	url_params := mux.Vars(req)
	following_id := url_params["user_id"]

	current_time := time.Now()

	follower_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	if follower_id.Value != following_id {
		_, err = ctrl.Db.Exec("insert into user_follows (follower_id, following_id, created_at) values ($1, $2, $3)", follower_id.Value, following_id, current_time)

		if err != nil {
			log.Println(err)
		}
	}

}

func (ctrl *FollowController) Delete(res http.ResponseWriter, req *http.Request) {
	follow := Models.Follow{}
	url_params := mux.Vars(req)
	following_id := url_params["user_id"]

	follower_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, follower_id, following_id, created_at from user_follows where follower_id = $1 and following_id = $2", follower_id.Value, following_id).Scan(&follow.Id, &follow.Follower_id, &follow.Following_id, &follow.Created_at)

	if err == sql.ErrNoRows {
		log.Printf("\nno follow record found:  ")
		log.Println(err)
	} else if err != nil {
		log.Println(err)
	} else {
		_, err = ctrl.Db.Exec("delete from user_follows where follower_id = $1 AND following_id = $2", follower_id.Value, following_id)

		if err != nil {
			log.Println(err)
		}
	}
}
