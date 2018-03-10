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
	gmux.HandleFunc("/follow-user/{user_id}", create_user_follow).Methods("GET")
	gmux.HandleFunc("/unfollow-user/{user_id}", delete_user_follow).Methods("GET")

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
		"Title":      "Home",
		"BodyHeader": "Welcome to the Starting Block",
		"Paragraph":  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.",
	}

	user_id, err := req.Cookie("session_uid")

	if err == nil {
		pageData["isUserLoggedIn"] = true

		if foundTweets, userTweets := getTweets(user_id.Value); foundTweets == true {
			pageData["foundTweets"] = true
			pageData["Tweets"] = userTweets
		} else {
			pageData["foundTweets"] = false
		}

	} else {
		pageData["isUserLoggedIn"] = false
	}

	fmt.Printf("\n")
	fmt.Println(pageData)
	fmt.Printf("\n")

	tpl.ExecuteTemplate(res, "index.html", pageData)
}

func getTweets(user_id string) (bool, []metaTweet) {
	fmt.Printf("\nGetting Tweets for user $s :o\n", user_id)

	var foundTweets = true
	var tweets []metaTweet

	rows, err := Db.Query("select distinct (t.id), t.user_id, u.name, u.username, t.msg, t.created_at from tweets t inner join users u on t.user_id = u.id inner join user_follows f on t.user_id = f.following_id where f.follower_id = $1 or t.user_id = $1 order by t.created_at desc", user_id)
	// old sql statement: select t.id, t.user_id, msg, t.created_at, count(*) from user_follows f left join tweets t on f.follower_id = $1 and f.following_id = t.user_id or t.user_id = $1 group by t.id order by t.created_at desc

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		tweet := metaTweet{}

		err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Name, &tweet.Username, &tweet.Message, &tweet.Created_at)

		switch {
		case err == sql.ErrNoRows:
			foundTweets = false
			log.Printf("\nNo tweets found.\n")

			return foundTweets, tweets
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
