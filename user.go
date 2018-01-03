package main

import (
	"fmt"
	// "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	// "text/template"
	// "time"
)

// Create a struct for querying User information
// Added bson tags to allow mgo to query mongoDB
type User struct {
	Name  string `bson:"name,omitempty"`
	Email string `bson:"email,omitempty"`
	Hash  string `bson:"hash,omitempty"`
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", r.URL.Path)

	r.ParseForm()
	name := r.PostFormValue("name")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	confirm_password := r.PostFormValue("confirm_password")
	new_user := User{}

	// if 'password' && 'confirm_password' have the same value
	if password == confirm_password {
		// Generate a hash from the submitted password with a cost of 10
		hashPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

		if err != nil {
			log.Println(err)
		}

		// Create the User struct were going to add into database
		new_user = User{name, email, string(hashPass)}

		/* if the hash we created and we take the same password and make it
		to a hash had the same original value, this should be true */
		if bcrypt.CompareHashAndPassword(hashPass, []byte(password)) == nil {
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
			// Insert the user struct, previously created, into the database
			col.Insert(new_user)
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}
