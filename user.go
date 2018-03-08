package main

import (
	// "database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

// Create a struct for querying User information
type User struct {
	Id         int
	Name       string
	Email      string
	Username   string
	Hash       string
	Created_at time.Time
	Updated_at time.Time
}

type User_meta struct {
	Id          int
	User_id     int
	Description string
	Url         string
	Created_at  time.Time
	Updated_at  time.Time
}

// GET
func user_register(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]string{
		"Title": "Sign Up",
	}

	// Execute the template
	tpl.ExecuteTemplate(res, "register.html", pageData)
}

// POST
func register(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	var user_id int
	req.ParseForm()
	name := req.PostFormValue("name")
	email := req.PostFormValue("email")
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	confirm_password := req.PostFormValue("confirm_password")

	// if 'password' && 'confirm_password' have the same value
	if password == confirm_password {
		// Generate a hash from the submitted password with a cost of 10
		hashPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		if err != nil {
			panic(err)
		}

		// Get the current time
		currentTime := time.Now()
		// insert the user into the users table in postgres
		nErr := Db.QueryRow("insert into users (name, email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $5) returning id", name, email, username, hashPass, currentTime).Scan(&user_id)

		// Check of there is an error connecting to the database
		if nErr == nil {
			_, err = Db.Exec("insert into user_meta (user_id, created_at, updated_at) values ($1, $2, $2)", user_id, currentTime)

			if err != nil {
				panic(nErr)
			}
		} else {
			panic(nErr)
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 200)
}

// GET
func update_user(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	retUser := User{}
	retMeta := User_meta{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select name, email, username from users where id = $1", user_id.Value).Scan(&retUser.Name, &retUser.Email, &retUser.Username)

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select description, url from user_meta where user_id = $1", user_id.Value).Scan(&retMeta.Description, &retMeta.Url)

	if err != nil {
		panic(err)
	}

	// Create map to pass data to template
	pageData := map[string]string{
		"Title":       "Settings",
		"Name":        retUser.Name,
		"Email":       retUser.Email,
		"Username":    retUser.Username,
		"Description": retMeta.Description,
		"Url":         retMeta.Url,
	}

	if is_user_logged_in(req) {
		pageData["isUserLoggedIn"] = "true"
	} else {
		pageData["isUserLoggedIn"] = "false"
	}

	tpl.ExecuteTemplate(res, "user_settings.html", pageData)
}

// POST
func change_user_info(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	name := req.PostFormValue("name")
	email := req.PostFormValue("email")
	username := req.PostFormValue("username")

	user_id, err := req.Cookie("session_uid")

	// Get the current time
	currentTime := time.Now()
	// update the users table in postgress
	_, err = Db.Exec("update users set name = $2, email = $3, username = $4, updated_at = $5 where id = $1", user_id.Value, name, email, username, currentTime)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/settings", 200)
}

// POST
func change_user_meta(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	bio := req.PostFormValue("bio")
	url := req.PostFormValue("url")

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	current_time := time.Now()

	_, err = Db.Exec("update user_meta set description = $2, url = $3, updated_at = $4 where user_id = $1", user_id.Value, bio, url, current_time)

	if err != nil {
		panic(err)
	}

	fmt.Printf("\nRedirecting to the '/settings' path\n")
	http.Redirect(res, req, "/settings", 200)
}

// POST
func change_user_password(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	old_password := req.PostFormValue("old-password")
	new_password := req.PostFormValue("new-password")
	confirm_new_password := req.PostFormValue("confirm-new-password")
	retUser := User{}

	user_id, err := req.Cookie("session_uid")

	err = Db.QueryRow("select id, password from users where id = $1;", user_id.Value).Scan(&retUser.Id, &retUser.Hash)

	if err != nil {
		panic(err)
	}

	// Compare the user hash and old_password
	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(old_password))

	// if ('pwd_match' doesnt have an error) and ('new_password' && 'confirm_password' have the same value)
	if (pwd_match == nil) && (new_password == confirm_new_password) {
		// Generate a hash from the submitted password with a cost of 10
		hashPass, err := bcrypt.GenerateFromPassword([]byte(new_password), 10)

		if err != nil {
			panic(err)
		}

		// Get the current time
		currentTime := time.Now()
		// update the user into the users table in postgres
		_, err = Db.Exec("update users set password = $2, updated_at = $3 where id = $1", user_id.Value, hashPass, currentTime)

		// Check of there is an error connecting to the database
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Passwords dont match in '/change_user_passwords'")
	}

	fmt.Printf("\nRedirecting to the '/settings' path\n")
	http.Redirect(res, req, "/settings", 200)
}

// POST
func delete_user(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	password := req.PostFormValue("password")
	retUser := User{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select id, password from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Hash)

	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	if pwd_match == nil {
		_, err := Db.Exec("delete from users where id = $1", retUser.Id)

		if err != nil {
			fmt.Println("Something went wrong. The user failed to be deleted:\n")
			fmt.Println(err)
		}
	}

	fmt.Printf("\nRedirecting to the '/settings' path\n")
	http.Redirect(res, req, "/logout", 200)
}

func show_user_profile(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	retUser := User{}
	retMeta := User_meta{}
	url_params := mux.Vars(req)
	user_id := url_params["user_id"]

	err := Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where id = $1", user_id).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select description, url from user_meta where user_id = $1", user_id).Scan(&retMeta.Description, &retMeta.Url)

	if err != nil {
		panic(err)
	}

	pageData := map[string]interface{}{
		"LoggedInUID": string(user_id),
		"ProfileUID":  string(retUser.Id),
		"Title":       retUser.Username + " | Profile",
		"Name":        retUser.Name,
		"Email":       retUser.Email,
		"Username":    retUser.Username,
		"Description": retMeta.Description,
		"Url":         retMeta.Url,
		"Tweets":      getUserTweets(user_id),
	}

	if err == nil && is_user_logged_in(req) {
		pageData["isUserLoggedIn"] = true
	} else {
		pageData["isUserLoggedIn"] = false
	}

	tpl.ExecuteTemplate(res, "profile.html", pageData)
}

func getUserTweets(user_id string) []Tweet {
	fmt.Printf("\nGetting Tweets :o\n")

	var tweets []Tweet

	// get user follow relations and use that to find all the tweets of the users the current logged in user follows, use 'group by' and 'count(*)' to do duplicate checking, and then order by the time created
	rows, err := Db.Query("select id, user_id, msg, created_at from tweets where user_id = $1 order by created_at desc limit 15", user_id)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		tweet := Tweet{}

		if err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Message, &tweet.Created_at); err != nil {
			log.Fatal(err)
		}

		tweets = append(tweets, tweet)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tweets
}
