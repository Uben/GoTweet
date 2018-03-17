package controllers

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gowebapp/helpers"
	"gowebapp/models"
	"log"
	"net/http"
	"time"
)

type UserController struct {
	Db *sql.DB
}

func NewUserController(DBCon *sql.DB) *UserController {
	return &UserController{DBCon}
}

func (ctrl *UserController) New(res http.ResponseWriter, req *http.Request) {
	pageData := map[string]interface{}{
		"Title":          "Sign Up",
		"isUserLoggedIn": false,
	}

	tpl.ExecuteTemplate(res, "register.html", pageData)
}

func (ctrl *UserController) Create(res http.ResponseWriter, req *http.Request) {
	var nil_check string
	current_time := time.Now()

	req.ParseForm()
	form := [5]string{req.PostFormValue("name"), req.PostFormValue("email"), req.PostFormValue("username"), req.PostFormValue("password"), req.PostFormValue("confirm_password")}

	if ((form[0] != nil_check) && (form[1] != nil_check) && (form[2] != nil_check) && (form[3] != nil_check) && (form[4] != nil_check)) && (form[3] == form[4]) {
		var user_id int
		hashPass, err := bcrypt.GenerateFromPassword([]byte(form[3]), 10)

		if err != nil {
			log.Println(err)
		}

		err = ctrl.Db.QueryRow("insert into users (name, email, username, password, created_at, updated_at) values ($1, $2, $3, $4, $5, $5) returning id", form[0], form[1], form[2], hashPass, current_time).Scan(&user_id)

		if err != nil {
			log.Println(err)
		}

		_, err = ctrl.Db.Exec("insert into user_meta (user_id, created_at, updated_at) values ($1, $2, $2)", user_id, current_time)

		if err != nil {
			log.Println(err)
		}
	}

	http.Redirect(res, req, "/", 302)
}

func (ctrl *UserController) Show(res http.ResponseWriter, req *http.Request) {
	retUser := Models.User{}
	retMeta := Models.UserMeta{}

	url_params := mux.Vars(req)
	user_id := url_params["user_id"]

	pageData := map[string]interface{}{
		"Title": "User Profile",
	}

	user_logged_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, name, email, username, password, created_at, updated_at from users where id = $1", user_id).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at, &retUser.Updated_at)

	switch {
	case err == sql.ErrNoRows:
		log.Println("No User record found")

		pageData = map[string]interface{}{
			"isUserValid":    false,
			"isUserLoggedIn": helpers.IsUserLoggedIn(req, ctrl.Db),
			"Title":          "User Not Found",
		}
	case err != nil:
		log.Fatal(err)

	default:
		log.Println("\nFound %s's info.\n", retUser.Name)

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

		if err == sql.ErrNoRows {
			log.Println("\nNo User meta record found\n")
		} else if err != nil {
			log.Fatal(err)
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
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select name, email, username from users where id = $1", user_id.Value).Scan(&retUser.Name, &retUser.Email, &retUser.Username)

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select description, url from user_meta where user_id = $1", user_id.Value).Scan(&retMeta.Description, &retMeta.Url)

	if err == sql.ErrNoRows {
		log.Println("\nNo User meta record found\n")
	} else if err != nil {
		log.Fatal(err)
	}

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
	current_time := time.Now()

	req.ParseForm()
	form := [3]string{req.PostFormValue("name"), req.PostFormValue("email"), req.PostFormValue("username")}

	if (form[0] != nil_check) && (form[1] != nil_check) && (form[2] != nil_check) {

		user_id, err := req.Cookie("session_uid")

		if err != nil {
			log.Println(err)
		}

		_, err = ctrl.Db.Exec("update users set name = $2, email = $3, username = $4, updated_at = $5 where id = $1", user_id.Value, form[0], form[1], form[2], current_time)

		if err != nil {
			log.Println(err)
		}

	} else {
		log.Println("\nNIL values were found. User NOT updated.\n")
	}

	http.Redirect(res, req, "/settings", 302)
}

func (ctrl *UserController) UpdateMeta(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	form := [2]string{req.PostFormValue("bio"), req.PostFormValue("url")}
	current_time := time.Now()

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	_, err = ctrl.Db.Exec("update user_meta set description = $2, url = $3, updated_at = $4 where user_id = $1", user_id.Value, form[0], form[1], current_time)

	if err != nil {
		log.Println(err)
	}

	http.Redirect(res, req, "/settings", 302)
}

func (ctrl *UserController) UpdatePassword(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	form := [3]string{req.PostFormValue("old-password"), req.PostFormValue("new-password"), req.PostFormValue("confirm-new-password")}

	var nil_check string
	retUser := Models.User{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	if (form[0] != nil_check) && (form[1] != nil_check) && (form[2] != nil_check) {

		err = ctrl.Db.QueryRow("select id, password from users where id = $1;", user_id.Value).Scan(&retUser.Id, &retUser.Hash)

		if err != nil {
			log.Println(err)
		}

		pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(form[0]))

		if (pwd_match == nil) && (form[1] == form[2]) {
			hashPass, err := bcrypt.GenerateFromPassword([]byte(form[1]), 10)
			current_time := time.Now()

			if err != nil {
				log.Println(err)
			}

			_, err = ctrl.Db.Exec("update users set password = $2, updated_at = $3 where id = $1", user_id.Value, hashPass, current_time)

			if err != nil {
				log.Println(err)
			}

		} else {
			log.Println("Passwords dont match in '/change_user_passwords'")
		}

	} else {
		log.Println("\nNIL values were found. User NOT updated.\n")
	}

	http.Redirect(res, req, "/settings", 302)
}

func (ctrl *UserController) Delete(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	password := req.PostFormValue("password")
	retUser := Models.User{}

	user_id, err := req.Cookie("session_uid")

	if err != nil {
		log.Println(err)
	}

	err = ctrl.Db.QueryRow("select id, password from users where id = $1", user_id.Value).Scan(&retUser.Id, &retUser.Hash)

	if err != nil {
		log.Println(err)
	}

	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	if pwd_match == nil {
		_, err := ctrl.Db.Exec("delete from users where id = $1", retUser.Id)

		if err != nil {
			log.Println(err)
		}

		_, err = ctrl.Db.Exec("delete from user_meta where user_id = $1", retUser.Id)

		if err != nil {
			log.Println(err)
		}
	}

	http.Redirect(res, req, "/logout", 302)
}
