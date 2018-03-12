package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gowebapp/models"
	"log"
	"net/http"
	"time"
)

// GET
func user_register(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]interface{}{
		"Title":          "Sign Up",
		"isUserLoggedIn": false,
	}

	// Execute the template
	tpl.ExecuteTemplate(res, "register.html", pageData)
}

// POST
func register(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	var nil_check string
	var user_id int
	req.ParseForm()
	name := req.PostFormValue("name")
	email := req.PostFormValue("email")
	username := req.PostFormValue("username")
	password := req.PostFormValue("password")
	confirm_password := req.PostFormValue("confirm_password")

	// if any of the submitted values weren't set, skip registering the user
	if (name != nil_check) && (email != nil_check) && (username != nil_check) && (password != nil_check) && (confirm_password != nil_check) {

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
					panic(err)
				}
			} else {
				panic(nErr)
			}

		}

	} else {
		fmt.Printf("\nnil values were found. User NOT created.\n")
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 200)
}

// GET
func update_user(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	retUser := Models.User{}
	retMeta := Models.UserMeta{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select name, email, username from users where id = $1", user_id.Value).Scan(&retUser.Name, &retUser.Email, &retUser.Username)

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select description, url from user_meta where user_id = $1", user_id.Value).Scan(&retMeta.Description, &retMeta.Url)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("\nNo usermeta found.\n")

	case err != nil:
		log.Fatal(err)

	default:
		fmt.Printf("\nFound %s's meta info.", retUser.Name)
	}

	// Create map to pass data to template
	pageData := map[string]interface{}{
		"Title":          "Settings",
		"Name":           retUser.Name,
		"Email":          retUser.Email,
		"Username":       retUser.Username,
		"Description":    retMeta.Description.String,
		"Url":            retMeta.Url.String,
		"isUserLoggedIn": is_user_logged_in(req),
	}

	tpl.ExecuteTemplate(res, "user_settings.html", pageData)
}

// POST
func change_user_info(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	var nil_check string
	req.ParseForm()
	name := req.PostFormValue("name")
	email := req.PostFormValue("email")
	username := req.PostFormValue("username")

	// if any of the submitted values weren't set, skip registering the user
	if (name != nil_check) && (email != nil_check) && (username != nil_check) {

		user_id, err := req.Cookie("session_uid")

		// Get the current time
		currentTime := time.Now()
		// update the users table in postgress
		_, err = Db.Exec("update users set name = $2, email = $3, username = $4, updated_at = $5 where id = $1", user_id.Value, name, email, username, currentTime)

		if err != nil {
			panic(err)
		}

	} else {

		fmt.Printf("\nNIL values were found. User NOT updated.\n")
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

	var nil_check string
	req.ParseForm()
	old_password := req.PostFormValue("old-password")
	new_password := req.PostFormValue("new-password")
	confirm_new_password := req.PostFormValue("confirm-new-password")
	retUser := Models.User{}

	user_id, err := req.Cookie("session_uid")

	if (old_password != nil_check) && (new_password != nil_check) && (confirm_new_password != nil_check) {

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
	} else {

		fmt.Printf("\nNIL values were found. User NOT updated.\n")
	}

	fmt.Printf("\nRedirecting to the '/settings' path\n")
	http.Redirect(res, req, "/settings", 200)
}

// POST
func delete_user(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	password := req.PostFormValue("password")
	retUser := Models.User{}

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

		_, err = Db.Exec("delete from user_meta where user_id = $1", retUser.Id)

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

	retUser := Models.User{}
	retMeta := Models.UserMeta{}
	url_params := mux.Vars(req)
	user_id := url_params["user_id"]

	user_logged_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	pageData := map[string]interface{}{
		"Title": "User Profile",
	}

	err = Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where id = $1", user_id).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No User record found")

		pageData = map[string]interface{}{
			"isUserValid":    false,
			"isUserLoggedIn": is_user_logged_in(req),
			"Title":          "User Not Found",
		}
	case err != nil:
		log.Fatal(err)

	default:
		fmt.Printf("\nFound %s's info.\n", retUser.Name)

		pageData = map[string]interface{}{
			"isUserValid":    true,
			"Title":          retUser.Username + " | Profile",
			"LoggedInUID":    string(user_logged_id.Value),
			"ProfileUID":     string(user_id),
			"Name":           retUser.Name,
			"Email":          retUser.Email,
			"Username":       retUser.Username,
			"isUserLoggedIn": is_user_logged_in(req),
		}

		err = Db.QueryRow("select description, url from user_meta where user_id = $1", user_id).Scan(&retMeta.Description, &retMeta.Url)

		switch {
		case err == sql.ErrNoRows:
			log.Printf("\nNo User meta record found\n")
		case err != nil:
			log.Fatal(err)
		default:
			fmt.Printf("\nFound %s's meta info.\n", retUser.Name)
		}

		if foundTweets, userTweets := getUserTweets(user_id); foundTweets == true {
			pageData["foundTweets"] = true
			pageData["Tweets"] = userTweets
		} else {
			pageData["foundTweets"] = false
		}

		if retMeta.Url.Valid == true {
			pageData["UrlSet"] = true
			pageData["Url"] = retMeta.Url.String
		} else {
			pageData["UrlSet"] = false
		}

		if retMeta.Description.Valid == true {
			pageData["DescriptionSet"] = true
			pageData["Description"] = retMeta.Url.String
		} else {
			pageData["DescriptionSet"] = false
		}
	}

	tpl.ExecuteTemplate(res, "profile.html", pageData)
}

func getUserTweets(user_id string) (bool, []Models.MetaTweet) {
	fmt.Printf("\nGetting Tweets :o\n")

	var foundTweets = true
	var tweets []Models.MetaTweet

	// get user follow relations and use that to find all the tweets of the users the current logged in user follows, use 'group by' and 'count(*)' to do duplicate checking, and then order by the time created
	rows, err := Db.Query("select distinct (t.id), t.user_id, u.name, u.username, t.msg, t.created_at from tweets t inner join users u on t.user_id = u.id where t.user_id = $1 order by t.created_at desc", user_id)
	// select id, user_id, msg, created_at from tweets where user_id = $1 order by created_at desc limit 15

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		tweet := Models.MetaTweet{}

		err := rows.Scan(&tweet.Id, &tweet.User_id, &tweet.Name, &tweet.Username, &tweet.Message, &tweet.Created_at)

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
