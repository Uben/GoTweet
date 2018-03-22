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

type TweetController struct {
	Db *sql.DB
}

func NewTweetController(DBCon *sql.DB) *TweetController {
	return &TweetController{DBCon}
}

func (ctrl *TweetController) Create(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	tweet_text := req.PostFormValue("tweet")

	retUser := Models.User{}
	current_time := time.Now()
	is_retweet := false

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

	if err != nil {
		log.Println(err)
	}

	_, err = ctrl.Db.Exec("insert into tweets (user_id, msg, name, username, is_retweet, origin_user_id, origin_name, origin_username, created_at) values ($1, $2, $3, $4, $5, $1, $3, $4, $6)", retUser.Id, tweet_text, retUser.Name, retUser.Username, is_retweet, current_time)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(res, req, "/", 302)
}

func (ctrl *TweetController) Delete(res http.ResponseWriter, req *http.Request) {
	session := Models.Session{}
	tweet := Models.Tweet{}

	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]

	session_id, err := req.Cookie("session")

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, user_id, token, created_at from sessions where token = $1", session_id.Value).Scan(&session.Id, &session.User_id, &session.Token, &session.Created_at)

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, user_id, msg, name, username, favorite_count, retweet_count, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at from tweets where id = $1", tweet_id).Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.FCount, &tweet.RCount, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

	if err == sql.ErrNoRows {
		log.Println(err)

	} else if err != nil {
		log.Println(err)

	} else if tweet.User_id == session.User_id {
		if tweet.Is_retweet {
			_, err := ctrl.Db.Exec("update tweets set retweet_count = retweet_count - 1 where id = $1 or origin_tweet_id = $1", tweet.Otweet_id)

			if err != nil {
				log.Println(err)
			}

			user_id, err := req.Cookie("session_uid")
			if err != nil {
				log.Println(err)
			}

			_, err = ctrl.Db.Exec("delete from retweets where user_id = $1 and tweet_id = $2", user_id.Value, tweet.Otweet_id)

			if err != nil {
				log.Println(err)
			}

			_, err = ctrl.Db.Exec("delete from tweets where id = $1", tweet_id)

			if err != nil {
				log.Println(err)
			}

		} else {
			_, err := ctrl.Db.Exec("delete from favorites where tweet_id = $1", tweet_id)

			if err != nil {
				log.Println(err)
			}

			_, err = ctrl.Db.Exec("delete from retweets where tweet_id = $1", tweet_id)

			if err != nil {
				log.Println(err)
			}

			_, err = ctrl.Db.Exec("delete from tweets where id = $1 or origin_tweet_id = $1", tweet_id)

			if err != nil {
				log.Println(err)
			}
		}
	}

	http.Redirect(res, req, "/", 302)
}

func (ctrl *TweetController) Retweet(res http.ResponseWriter, req *http.Request) {
	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]

	tweet := Models.Tweet{}
	retUser := Models.User{}
	is_retweet := true
	current_time := time.Now()

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, user_id, msg, name, username, favorite_count, retweet_count, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at from tweets where id = $1 limit 1", tweet_id).Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.FCount, &tweet.RCount, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, name, username from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Name, &retUser.Username)

	if err != nil {
		log.Println(err)
	}

	if tweet.Is_retweet {
		_, err := ctrl.Db.Exec("insert into tweets (user_id, msg, name, username, favorite_count, retweet_count, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", retUser.Id, tweet.Message, retUser.Name, retUser.Username, tweet.FCount, tweet.RCount, is_retweet, tweet.Otweet_id, tweet.Ouser_id, tweet.Oname, tweet.Ousername, current_time)

		if err != nil {
			log.Println(err)

		} else if err == nil {
			_, err = ctrl.Db.Exec("insert into retweets (user_id, tweet_id, created_at) values ($1, $2, $3)", user_id.Value, tweet.Otweet_id, current_time)

			if err != nil {
				log.Println(err)
			}
		}

	} else {
		_, err := ctrl.Db.Exec("insert into tweets (user_id, msg, name, username, favorite_count, retweet_count, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", retUser.Id, tweet.Message, retUser.Name, retUser.Username, tweet.FCount, tweet.RCount, is_retweet, tweet.Id, tweet.User_id, tweet.Name, tweet.Username, current_time)

		if err != nil {
			log.Println(err)

		} else if err == nil {
			_, err := ctrl.Db.Exec("insert into retweets (user_id, tweet_id, created_at) values ($1, $2, $3)", user_id.Value, tweet_id, current_time)

			if err != nil {
				log.Println(err)
			}
		}
	}

	_, err = ctrl.Db.Exec("update tweets set retweet_count = retweet_count + 1 where id = $1 or origin_tweet_id = $1", tweet_id)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(res, req, "/", 302)
}

func (ctrl *TweetController) Favorite(res http.ResponseWriter, req *http.Request) {
	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]

	current_time := time.Now()

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	_, err = ctrl.Db.Exec("insert into favorites (user_id, tweet_id, created_at) values ($1, $2, $3)", user_id.Value, tweet_id, current_time)

	if err != nil {
		log.Println(err)
	}

	_, err = ctrl.Db.Exec("update tweets set favorite_count = favorite_count + 1 where id = $1 or origin_tweet_id = $1", tweet_id)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(res, req, "/", 302)
}

func (ctrl *TweetController) Unfavorite(res http.ResponseWriter, req *http.Request) {
	url_params := mux.Vars(req)
	tweet_id := url_params["tweet_id"]

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	_, err = ctrl.Db.Exec("delete from favorites where user_id = $1 and tweet_id = $2", user_id.Value, tweet_id)

	if err != nil {
		log.Println(err)
	}

	_, err = ctrl.Db.Exec("update tweets set favorite_count = favorite_count - 1 where id = $1 or origin_tweet_id = $1", tweet_id)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(res, req, "/", 302)
}
