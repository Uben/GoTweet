package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template
var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=devtest dbname=gowebapp password=password sslmode=disable")

	if err != nil {
		panic(err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	gmux := mux.NewRouter()

	gmux.HandleFunc("/register", user_register).Methods("GET")
	gmux.HandleFunc("/register", register).Methods("POST")

	gmux.HandleFunc("/login", user_login).Methods("GET")
	gmux.HandleFunc("/login", login).Methods("POST")
	gmux.HandleFunc("/logout", isAuth(logout)).Methods("GET")

	gmux.HandleFunc("/settings", update_user).Methods("GET")
	gmux.HandleFunc("/update-user-info", change_user_info).Methods("POST")
	gmux.HandleFunc("/update-user-meta", change_user_meta).Methods("POST")
	gmux.HandleFunc("/update-user-password", change_user_password).Methods("POST")

	gmux.HandleFunc("/create-tweet", tweet_create).Methods("POST")
	gmux.HandleFunc("/delete-tweet/{tweet_id}", tweet_delete).Methods("POST")

	gmux.HandleFunc("/profile/{user_id}", show_user_profile).Methods("GET")
	gmux.HandleFunc("/follow-user/{user_id}", create_user_follow).Methods("POST")
	gmux.HandleFunc("/unfollow-user/{user_id}", delete_user_follow).Methods("POST")

	gmux.HandleFunc("/favicon.ico", handlerIcon).Methods("GET")
	gmux.HandleFunc("/", home).Methods("GET")

	fmt.Printf("About to listen on port :3000. Go to https://127.0.0.1:3000/ (localhost)\n")
	log.Fatal(http.ListenAndServe(":3000", gmux))
}

func handlerIcon(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)
}

func home(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]interface{}{
		"Title":      "Home | Base Golang Web App",
		"BodyHeader": "Welcome to the Starting Block",
		"Paragraph":  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.",
	}

	user_id, err := req.Cookie("session_uid")

	if err == nil && is_user_logged_in(req) {
		pageData["isUserLoggedIn"] = "true"
		pageData["Tweets"] = getTweets(user_id.Value)
	} else {
		pageData["isUserLoggedIn"] = "false"
	}

	fmt.Println(pageData)

	// Check if the path is exactly "/" else its a 404 error
	if req.URL.Path != "/" {
		tpl.ExecuteTemplate(res, "404.html", pageData)

		// Else Execute the index template with the 'pageData' data
	} else {
		tpl.ExecuteTemplate(res, "index.html", pageData)
	}
}

func getTweets(user_id string) []Tweet {
	fmt.Printf("\nGetting Tweets :o\n")

	var tweets []Tweet

	// get user follow relations and use that to find all the tweets of the users the current logged in user follows, use 'group by' and 'count(*)' to do duplicate checking, and then order by the time created
	rows, err := Db.Query("select t.id, t.user_id, msg, t.created_at, count(*) from user_follows f left join tweets t on f.follower_id = $1 and f.following_id = t.user_id or t.user_id = $1 group by t.id order by t.created_at desc", user_id)

	if err == sql.ErrNoRows {
		rows.Close()

	} else if err != nil {
		log.Fatal(err)

	} else {
		defer rows.Close()

		for rows.Next() {
			tweet := Tweet{}
			var throwaway int

			if err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Created_at, &throwaway); err != nil {
				// log.Fatal(err)
			}

			tweets = append(tweets, tweet)
		}

		if err := rows.Err(); err != nil {
			// log.Fatal(err)
		}
	}

	return tweets
}
