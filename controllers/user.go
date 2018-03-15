package controllers

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gowebapp/helpers"
	"gowebapp/models"
	"log"
	"net/http"
	"text/template"
	"time"
)

type UserController struct {
	Db *sql.DB
}

func NewUserController(DBCon *sql.DB) *UserController {
	return &UserController{DBCon}
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func (ctrl *UserController) New(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]interface{}{
		"Title":          "Sign Up",
		"isUserLoggedIn": false,
	}

	// Execute the template
	tpl.ExecuteTemplate(res, "register.html", pageData)
}

func (ctrl *UserController) Create(res http.ResponseWriter, req *http.Request) {

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
			current_time := time.Now()
			// insert the user into the users table in postgres
			err = ctrl.Db.QueryRow("insert into users (name, email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $5) returning id", name, email, username, hashPass, current_time).Scan(&user_id)

			// Check of there is an error connecting to the database
			if err == nil {
				_, err = ctrl.Db.Exec("insert into user_meta (user_id, created_at, updated_at) values ($1, $2, $2)", user_id, current_time)

				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}

		}

	} else {
		fmt.Printf("\nnil values were found. User NOT created.\n")
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/", 200)
}

func (ctrl *UserController) Show(res http.ResponseWriter, req *http.Request) {
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

	err = ctrl.Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where id = $1", user_id).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("No User record found")

		pageData = map[string]interface{}{
			"isUserValid":    false,
			"isUserLoggedIn": helpers.IsUserLoggedIn(req, ctrl.Db),
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
			"isUserLoggedIn": helpers.IsUserLoggedIn(req, ctrl.Db),
		}

		err = ctrl.Db.QueryRow("select description, url from user_meta where user_id = $1", user_id).Scan(&retMeta.Description, &retMeta.Url)

		switch {
		case err == sql.ErrNoRows:
			log.Printf("\nNo User meta record found\n")
		case err != nil:
			log.Fatal(err)
		default:
			fmt.Printf("\nFound %s's meta info.\n", retUser.Name)
		}

		if foundTweets, userTweets := helpers.GetUserTweets(user_id, ctrl.Db); foundTweets == true {
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

func (ctrl *UserController) Edit(res http.ResponseWriter, req *http.Request) {
	retUser := Models.User{}
	retMeta := Models.UserMeta{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	err = ctrl.Db.QueryRow("select name, email, username from users where id = $1", user_id.Value).Scan(&retUser.Name, &retUser.Email, &retUser.Username)

	if err != nil {
		panic(err)
	}

	err = ctrl.Db.QueryRow("select description, url from user_meta where user_id = $1", user_id.Value).Scan(&retMeta.Description, &retMeta.Url)

	if err != nil {
		panic(err)
	}

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
		"isUserLoggedIn": helpers.IsUserLoggedIn(req, ctrl.Db),
	}

	tpl.ExecuteTemplate(res, "user_settings.html", pageData)
}

func (ctrl *UserController) UpdateInfo(res http.ResponseWriter, req *http.Request) {
	var nil_check string
	req.ParseForm()
	name := req.PostFormValue("name")
	email := req.PostFormValue("email")
	username := req.PostFormValue("username")

	// if any of the submitted values weren't set, skip registering the user
	if (name != nil_check) && (email != nil_check) && (username != nil_check) {

		user_id, err := req.Cookie("session_uid")

		// Get the current time
		current_time := time.Now()
		// update the users table in postgress
		_, err = ctrl.Db.Exec("update users set name = $2, email = $3, username = $4, updated_at = $5 where id = $1", user_id.Value, name, email, username, current_time)

		if err != nil {
			panic(err)
		}

	} else {

		fmt.Printf("\nNIL values were found. User NOT updated.\n")
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(res, req, "/settings", 302)
}

func (ctrl *UserController) UpdateMeta(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	bio := req.PostFormValue("bio")
	url := req.PostFormValue("url")

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	current_time := time.Now()

	_, err = ctrl.Db.Exec("update user_meta set description = $2, url = $3, updated_at = $4 where user_id = $1", user_id.Value, bio, url, current_time)

	if err != nil {
		panic(err)
	}

	http.Redirect(res, req, "/settings", 302)
}

func (ctrl *UserController) UpdatePassword(res http.ResponseWriter, req *http.Request) {
	var nil_check string
	req.ParseForm()
	old_password := req.PostFormValue("old-password")
	new_password := req.PostFormValue("new-password")
	confirm_new_password := req.PostFormValue("confirm-new-password")
	retUser := Models.User{}

	user_id, err := req.Cookie("session_uid")

	if (old_password != nil_check) && (new_password != nil_check) && (confirm_new_password != nil_check) {

		err = ctrl.Db.QueryRow("select id, password from users where id = $1;", user_id.Value).Scan(&retUser.Id, &retUser.Hash)

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
			current_time := time.Now()
			// update the user into the users table in postgres
			_, err = ctrl.Db.Exec("update users set password = $2, updated_at = $3 where id = $1", user_id.Value, hashPass, current_time)

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

	http.Redirect(res, req, "/settings", 302)
}

func (ctrl *UserController) Delete(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	password := req.PostFormValue("password")
	retUser := Models.User{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		panic(err)
	}

	err = ctrl.Db.QueryRow("select id, password from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Hash)

	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	if pwd_match == nil {
		_, err := ctrl.Db.Exec("delete from users where id = $1", retUser.Id)

		if err != nil {
			fmt.Println("Something went wrong. The user failed to be deleted:\n")
			fmt.Println(err)
		}

		_, err = ctrl.Db.Exec("delete from user_meta where user_id = $1", retUser.Id)

		if err != nil {
			fmt.Println("Something went wrong. The user failed to be deleted:\n")
			fmt.Println(err)
		}
	}

	fmt.Printf("\nRedirecting to the '/settings' path\n")
	http.Redirect(res, req, "/logout", 302)
}
