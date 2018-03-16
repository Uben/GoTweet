package helpers

import (
	"database/sql"
	"fmt"
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

	// get user follow relations and use that to find all the tweets of the users the current logged in user follows, use 'group by' and 'count(*)' to do duplicate checking, and then order by the time created
	rows, err := Db.Query("select id, user_id, msg, name, username, is_retweet, origin_tweet_id, origin_user_id, origin_name, origin_username, created_at from tweets where user_id = $1 order by created_at desc", user_id)
	// select id, user_id, msg, created_at from tweets where user_id = $1 order by created_at desc limit 15
	// select distinct (t.id), t.user_id, u.name, u.username, t.msg, t.is_retweet, origin_tweet_id, origin_user_id, t.created_at from tweets t inner join users u on t.user_id = u.id where t.user_id = $1 order by t.created_at desc

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		tweet := Models.Tweet{}

		err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Name, &tweet.Username, &tweet.Is_retweet, &tweet.Otweet_id, &tweet.Ouser_id, &tweet.Oname, &tweet.Ousername, &tweet.Created_at)

		switch {
		case err == sql.ErrNoRows:
			foundTweets = false
			log.Printf("No user with that ID.")
		case err != nil:
			log.Fatal(err)
		default:
			fmt.Printf("\nAdded a tweet.")
			tweets = append(tweets, tweet)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return foundTweets, tweets
}

func IsUserLoggedIn(r *http.Request, Db *sql.DB) bool {
	fmt.Printf("\nChecking if the user is logged in.\n")

	// Get the "session" cookie
	session_cookie, err := r.Cookie("session")
	// Create an empty session struct
	retSession := Models.Session{}

	// If there isnt an error && the value of 'session_cookie' isnt equal to ""
	if (err == nil) && (session_cookie.Value != "") {

		// Find a document with a 'Token' value that is equal to the session cookie value
		err := Db.QueryRow("select id, user_id, token, created_at from sessions where token = $1", session_cookie.Value).Scan(&retSession.Id, &retSession.User_id, &retSession.Token, &retSession.Created_at)

		// If there is no error getting the data && Check if the value of 'uuid' in the found document is equal to the 'Value' in 'session_cookie'
		if (err == nil) && (retSession.Token == session_cookie.Value) {
			fmt.Printf("\nUser is logged in with session: '%s'.", session_cookie.Value)
			return true
		} else {
			fmt.Println("\nUser is NOT logged in.\n")
			return false
		}

	} else {
		fmt.Println("\nUser is NOT logged in.\n")
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
		panic(err)
	}

	return true
}
