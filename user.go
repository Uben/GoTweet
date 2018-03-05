package main

import (
	// "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

// GET
func user_register(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]string{
		"title": "Sign Up",
	}

	// Execute the template
	tpl.ExecuteTemplate(res, "register.html", pageData)
}

// POST
func register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", r.URL.Path)

	r.ParseForm()
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	confirm_password := r.PostFormValue("confirm_password")

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
		_, nErr := Db.Exec("insert into users (name, email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $5)", name, email, username, hashPass, currentTime)

		// Check of there is an error connecting to the database
		if nErr != nil {
			panic(nErr)
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}
