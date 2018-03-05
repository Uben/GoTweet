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
func register(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", req.URL.Path)

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
		_, nErr := Db.Exec("insert into users (name, email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $5)", name, email, username, hashPass, currentTime)

		// Check of there is an error connecting to the database
		if nErr != nil {
			panic(nErr)
		}
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 200)
}

// GET
func update_user(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	user_id, err := req.Cookie("session_uid")
	retUser := User{}

	if err != nil {
		panic(err)
	}

	err = Db.QueryRow("select id, name, email, username from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username)

	// Create map to pass data to template
	pageData := map[string]string{
		"Title":    "Account Settings | Base Golang Web App",
		"Name":     retUser.Name,
		"Email":    retUser.Email,
		"Username": retUser.Username,
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
