package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gowebapp/controllers"
	"log"
	"net/http"
)

var Db *sql.DB

func init() {
	var err error

	if Db, err = sql.Open("postgres", "user=devtest dbname=gowebapp password=password sslmode=disable"); err != nil {
		panic(err)
	}
}

func main() {
	// Controllers
	gc := controllers.NewGeneralController(Db)
	uc := controllers.NewUserController(Db)
	sc := controllers.NewSessionController(Db)
	fc := controllers.NewFollowController(Db)
	tc := controllers.NewTweetController(Db)

	gmux := mux.NewRouter()

	// Registration Routes
	gmux.HandleFunc("/register", uc.New).Methods("GET")
	gmux.HandleFunc("/register", uc.Create).Methods("POST")
	gmux.HandleFunc("/delete-user", uc.Delete).Methods("POST")

	// Session Routes
	gmux.HandleFunc("/login", sc.Get).Methods("GET")
	gmux.HandleFunc("/login", sc.Create).Methods("POST")
	gmux.HandleFunc("/logout", sc.Delete).Methods("GET")

	// User Account Info Routes
	gmux.HandleFunc("/settings", uc.Edit).Methods("GET")
	gmux.HandleFunc("/update-user-info", uc.UpdateInfo).Methods("POST")
	gmux.HandleFunc("/update-user-meta", uc.UpdateMeta).Methods("POST")
	gmux.HandleFunc("/update-user-password", uc.UpdatePassword).Methods("POST")

	gmux.HandleFunc("/profile/{user_id}", uc.Show).Methods("GET")
	gmux.HandleFunc("/follow-user/{user_id}", fc.Create).Methods("GET")
	gmux.HandleFunc("/unfollow-user/{user_id}", fc.Delete).Methods("GET")

	/* Tweet Routes */
	gmux.HandleFunc("/create-tweet", tc.Create).Methods("POST")
	gmux.HandleFunc("/delete-tweet/{tweet_id}", tc.Delete).Methods("GET")
	gmux.HandleFunc("/create-retweet/{tweet_id}", tc.Retweet).Methods("GET")

	/* Tweet Favorite Routes */
	gmux.HandleFunc("/favorite/{tweet_id}", tc.Favorite).Methods("GET")
	gmux.HandleFunc("/unfavorite/{tweet_id}", tc.Unfavorite).Methods("GET")

	/* Base Routes */
	gmux.HandleFunc("/favicon.ico", gc.Favicon).Methods("GET")
	gmux.HandleFunc("/", gc.Home).Methods("GET")

	log.Printf("- About to listen on port :3000. Go to http://127.0.0.1:3000/ (http://localhost:3000/)\n")
	log.Fatal(http.ListenAndServe(":3000", gmux))
}
