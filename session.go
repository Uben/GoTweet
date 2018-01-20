package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// Create a struct for querying Session information
// Added bson tags to allow mgo to query mongoDB
type Session struct {
	Email string `bson:"email,omitempty"`
	Uuid  string `bson:"uuid,omitempty"`
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", r.URL.Path)

	r.ParseForm()
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	retUser := User{}

	// Create a connection to the database
	conn, err := mgo.Dial("mongodb://127.0.0.1:27017")
	// Check of there is an error connecting to the database
	if err != nil {
		panic(err)
	}
	// Defer closing the connection to the end of this function
	defer conn.Close()

	// Connect to the database "webdev" and access the "users" collection
	col := conn.DB("webdev").C("users")

	// Find one user with the email value of 'email' in the "users" collection
	err = col.Find(bson.M{"email": string(email)}).One(&retUser)

	if err != nil {
		panic(err)
	}

	// Compare the user hash and password
	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	// if the password match ^ (above) ^ doesnt return a error then continue
	if pwd_match == nil {
		fmt.Printf("\nCreating the session for the user '%s'.\n", retUser.Name)

		// // Get the "session" cookie
		// session_cookie, err := r.Cookie("session")

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

			// Connect to the database "webdev" and access the "sessions" collection
			sCol := conn.DB("webdev").C("sessions")

			// Insert the session struct in the 'session' collection in the database
			sCol.Insert(Session{email, user_uuid.String()})

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

	// Create a connection to the database
	conn, err := mgo.Dial("mongodb://127.0.0.1:27017")
	// Check of there is an error connecting to the database
	if err != nil {
		panic(err)
	}
	// Defer closing the connection to the end of this function
	defer conn.Close()

	// Connect to the database "webdev" and access the "sessions" collection
	sCol := conn.DB("webdev").C("sessions")

	// Get the "session" cookie
	session_cookie, err := r.Cookie("session")

	if err == nil {
		// Remove the document where the "uuid" field is the same as the users session UUID value
		sCol.Remove(bson.M{"uuid": session_cookie.Value})

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

	// Create a connection to the database
	conn, err := mgo.Dial("mongodb://127.0.0.1:27017")
	// Check of there is an error connecting to the database
	if err != nil {
		panic(err)
	}
	// Defer closing the connection to the end of this function
	defer conn.Close()

	// Connect to the database "webdev" and access the "sessions" collection
	sCol := conn.DB("webdev").C("sessions")

	// Get the "session" cookie
	session_cookie, err := r.Cookie("session")
	// Create an empty session struct
	retSession := Session{}

	// If there isnt an error && the value of 'session_cookie' isnt equal to ""
	if (err == nil) && (session_cookie.Value != "") {
		// Find a document with a 'Uuid' value that is equal to the session cookie
		err = sCol.Find(bson.M{"uuid": session_cookie.Value}).One(&retSession)

		// If there is no error getting the data
		if err == nil {
			// Check if the value of 'uuid' in the found document is equal to the 'Value' in 'session_cookie'
			if retSession.Uuid == session_cookie.Value {
				fmt.Println("\nSession: ", retSession)
				fmt.Printf("\nUser '%s' is logged in with session: '%s'.", retSession.Email, session_cookie.Value)
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
