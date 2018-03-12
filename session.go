package main

import (
	// "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gowebapp/models"
	"net/http"
	"strconv"
	"time"
)

// GET
func user_login(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]interface{}{
		"Title":          "Login",
		"isUserLoggedIn": false,
	}

	// Execute the template
	tpl.ExecuteTemplate(res, "login.html", pageData)
}

// POST
func login(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	req.ParseForm()
	email := req.PostFormValue("email")
	password := req.PostFormValue("password")
	retUser := Models.User{}

	// Get the user info and scan it into the user struct
	err := Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where email = $1 limit 1", email).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

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
			user_uuid, err := uuid.NewV4()

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

			// Set the "session" cookie values
			session_uid_cookie := &http.Cookie{
				Name:     "session_uid",
				Value:    strconv.Itoa(retUser.Id),
				HttpOnly: true,
				Path:     "/",
				Expires:  expire,
				MaxAge:   86400,
			}

			// Set the "session" cookie values
			session_username_cookie := &http.Cookie{
				Name:     "session_username",
				Value:    retUser.Username,
				HttpOnly: true,
				Path:     "/",
				Expires:  expire,
				MaxAge:   86400,
			}

			// Set the Cookies
			http.SetCookie(res, session_cookie)
			http.SetCookie(res, session_uid_cookie)
			http.SetCookie(res, session_username_cookie)

			currentTime := time.Now()
			_, err = Db.Exec("insert into sessions (user_id, token, created_at, updated_at) values ($1, $2, $3, $3)", &retUser.Id, user_uuid.String(), currentTime)

			if err != nil {
				panic(err)
			}

			fmt.Printf("\nUser: %s, has logged in with Session ID UUID: '%s'", retUser.Name, user_uuid)
		} else {
			// Get the "session" cookie
			session_cookie, err := req.Cookie("session")

			if err == nil {
				fmt.Printf("\nUser: %s, is ALREADY logged in with Session ID UUID: '%s'", retUser.Name, session_cookie.Value)
			}

		}

	} else {

		fmt.Printf("\nNAH, YOU TRIED IT. YAH FAILED. GO AGAIN.")

	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 303)
}

// GET
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

		// Set the "session" cookie values
		session_uid_cookie := &http.Cookie{
			Name:     "session_uid",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -10,
			Expires:  time.Now(),
		}

		// Set the "session" cookie values
		session_username_cookie := &http.Cookie{
			Name:     "session_username",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -10,
			Expires:  time.Now(),
		}

		// Set the Cookies
		http.SetCookie(w, session_cookie)
		http.SetCookie(w, session_uid_cookie)
		http.SetCookie(w, session_username_cookie)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}

func isAuth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check if user is logged in
		userAuthStatus := is_user_logged_in(req)

		// if the user is logged in, show the page
		if userAuthStatus == true {
			h.ServeHTTP(res, req)

			// else have them login
		} else {
			http.Redirect(res, req, "/", 200)
		}

	})
}

func is_user_logged_in(r *http.Request) bool {
	fmt.Printf("\nChecking if the user is logged in.\n")

	// Get the "session" cookie
	session_cookie, err := r.Cookie("session")
	// Create an empty session struct
	retSession := Models.Session{}

	// If there isnt an error && the value of 'session_cookie' isnt equal to ""
	if (err == nil) && (session_cookie.Value != "") {

		// Find a document with a 'Token' value that is equal to the session cookie value
		err := Db.QueryRow("select id, user_id, token, created_at, updated_at from sessions where token = $1", session_cookie.Value).Scan(&retSession.Id, &retSession.User_id, &retSession.Token, &retSession.Created_at, &retSession.Updated_at)

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
}
