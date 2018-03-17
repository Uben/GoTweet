package helpers

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gowebapp/models"
	"log"
	"net/http"
	"strconv"
)

var Db *sql.DB

func GetUserTweets(user_id string, Db *sql.DB) (bool, []Models.Tweet) {
	var foundTweets = true
	var tweets []Models.Tweet

	rows, err := Db.Query("select id, user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at from tweets where user_id = $1 order by created_at desc", user_id)

	if err != nil {
		log.Fatal(err)
		return false, tweets
	}

	defer rows.Close()

	for rows.Next() {
		tweet := Models.Tweet{}

		err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

		switch {
		case err == sql.ErrNoRows:
			foundTweets = false
		case err != nil:
			log.Fatal(err)
		default:
			tweets = append(tweets, tweet)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return foundTweets, tweets
}

func GetTweets(user_id string, Db *sql.DB) (bool, []Models.Tweet) {
	var foundTweets = true
	var tweets []Models.Tweet

	rows, err := Db.Query("select distinct (t.id), user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, t.created_at from tweets t inner join user_follows f on t.user_id = f.following_id where f.follower_id = $1 or t.user_id = $1 order by t.created_at desc", user_id)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		tweet := Models.Tweet{}

		err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

		switch {
		case err == sql.ErrNoRows:
			foundTweets = false
			return foundTweets, tweets
		case err != nil:
			log.Fatal(err)
		default:
			tweets = append(tweets, tweet)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return foundTweets, tweets
}

func IsAuth(handle http.HandlerFunc, Db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if userAuthStatus := IsUserLoggedIn(req, Db); userAuthStatus == true {
			handle.ServeHTTP(res, req)
		} else {
			http.Redirect(res, req, "/", 200)
		}
	})
}

func IsUserLoggedIn(r *http.Request, Db *sql.DB) bool {
	session_cookie, err := r.Cookie("session")
	retSession := Models.Session{}

	if (err == nil) && (session_cookie.Value != "") {
		err := Db.QueryRow("select id, user_id, token, created_at from sessions where token = $1", session_cookie.Value).Scan(&retSession.Id, &retSession.User_id, &retSession.Token, &retSession.Created_at)

		if (err == nil) && (retSession.Token == session_cookie.Value) {
			log.Printf("\nUser is logged in with session: '%s'.", session_cookie.Value)
			return true
		} else {
			log.Printf("\nUser is NOT logged in.\n")
			return false
		}

	} else {
		log.Printf("\nUser is NOT logged in.\n")
		return false
	}

	return false
}

func IsFavorite(tweet_id, user_id int, Db *sql.DB) bool {
	fave := Models.Favorite{}
	err := Db.QueryRow("select id, user_id, tweet_id, created_at from favorites where user_id = $1 and tweet_id = $2", strconv.Itoa(user_id), strconv.Itoa(tweet_id)).Scan(&fave.Id, &fave.User_id, &fave.Tweet_id, &fave.Created_at)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return true
}
