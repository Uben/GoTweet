package main

import (
	// "database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Create a struct for querying User information
type User struct {
	Id    int
	Name  string
	Email string
	Hash  string
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", r.URL.Path)

	r.ParseForm()
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	confirm_password := r.PostFormValue("confirm_password")

	// if 'password' && 'confirm_password' have the same value
	if password == confirm_password {
		// Generate a hash from the submitted password with a cost of 10
		hashPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		if err != nil {
			panic(err)
		}

		// insert the user into the users table in postgres
		_, nErr := Db.Exec("insert into users (name, email, password) values ($1, $2, $3)", name, email, hashPass)

		// Check of there is an error connecting to the database
		if nErr != nil {
			panic(nErr)
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}
