package main

import (
	// "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// Create a struct for querying Session information
// Added bson tags to allow mgo to query mongoDB
type Session struct {
	Id      int
	User_id int
	Token   string
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", r.URL.Path)

	r.ParseForm()
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	retUser := User{}

	// Get the user info and scan it into the user struct
	err := Db.QueryRow("select id, name, email, password from users where email = $1 limit 1", email).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Hash)

	if err != nil {
		panic(err)
	}

	// Compare the user hash and password
	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	// if the password match ^ (above) ^ doesnt return a error then continue
	if pwd_match == nil {
		fmt.Printf("\nCreating the session for the user '%s'.\n", retUser.Name)

		if err == nil {
			// Create a new UUID for the session
			user_uuid := uuid.NewV4()

			// Set cookie Expire date to one day from now
			expire := time.Now().AddDate(0, 0, 1)

			// Set the "session" cookie values
			session_cookie := &http.Cookie{
				Name:     "session",
				Value:    user_uuid.String(),
				HttpOnly: true,
				Path:     "/",
				Expires:  expire,
				MaxAge:   86400,
			}

			// Set the Cookie
			http.SetCookie(w, session_cookie)

			_, err := Db.Exec("insert into sessions (user_id, token) values ($1, $2)", &retUser.Id, user_uuid.String())

			if err != nil {
				panic(err)
			}

			fmt.Printf("\nUser: %s, has logged in with Session ID UUID: '%s'", retUser.Name, user_uuid)
		} else {
			// Get the "session" cookie
			session_cookie, err := r.Cookie("session")

			if err == nil {
				fmt.Printf("\nUser: %s, is ALREADY logged in with Session ID UUID: '%s'", retUser.Name, session_cookie.Value)
			}

		}

	} else {

		fmt.Printf("\nNAH, YOU TRIED IT. YAH FAILED. GO AGAIN.")

	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", r.URL.Path)

	// Get the "session" cookie
	session_cookie, err := r.Cookie("session")

	if err == nil {
		_, err := Db.Exec("delete from sessions where token = $1", session_cookie.Value)

		if err != nil {
			panic(err)
		}

		// Set the "session" cookie values
		session_cookie = &http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -10,
			Expires:  time.Now(),
		}
		// Set the Cookie
		http.SetCookie(w, session_cookie)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}

func is_user_logged_in(r *http.Request) (bool, int) {
	fmt.Printf("\nChecking if the user is logged in.\n")

	// Get the "session" cookie
	session_cookie, err := r.Cookie("session")
	// Create an empty session struct
	retSession := Session{}

	// If there isnt an error && the value of 'session_cookie' isnt equal to ""
	if (err == nil) && (session_cookie.Value != "") {

		// Find a document with a 'Token' value that is equal to the session cookie value
		err := Db.QueryRow("select id, user_id, token from sessions where token = $1", session_cookie.Value).Scan(&retSession.Id, &retSession.User_id, &retSession.Token)

		// If there is no error getting the data
		if err == nil {
			// Check if the value of 'uuid' in the found document is equal to the 'Value' in 'session_cookie'
			if retSession.Token == session_cookie.Value {
				fmt.Println("\nSession: ", retSession)
				fmt.Printf("\nUser is logged in with session: '%s'.", session_cookie.Value)
				return true, 1
			} else {
				fmt.Println("\nUser is NOT logged in.\n")
				return false, 0
			}
		} else {
			panic(err)
		}

	} else {
		fmt.Println("\nUser is NOT logged in.\n")
		return false, 0
	}
}
