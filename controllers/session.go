package controllers

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gowebapp/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

type SessionController struct {
	Db *sql.DB
}

func NewSessionController(DBCon *sql.DB) *SessionController {
	return &SessionController{DBCon}
}

func (ctrl *SessionController) Get(res http.ResponseWriter, req *http.Request) {
	pageData := map[string]interface{}{
		"Title":          "Login",
		"isUserLoggedIn": false,
	}

	tpl.ExecuteTemplate(res, "login.html", pageData)
}

func (ctrl *SessionController) Create(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	email := req.PostFormValue("email")
	password := req.PostFormValue("password")
	retUser := Models.User{}

	err := ctrl.Db.QueryRow("select id, name, email, username, password, created_at from users where email = $1 limit 1", email).Scan(&retUser.Id, &retUser.Name, &retUser.Email, &retUser.Username, &retUser.Hash, &retUser.Created_at)

	if err != nil {
		log.Println(err)
	}

	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	if pwd_match == nil {
		log.Println("\nCreating the session for the user '%s'.\n", retUser.Name)

		if err == nil {
			current_time := time.Now()
			user_uuid, err := uuid.NewV4()
			expire := time.Now().AddDate(0, 0, 1)

			session_cookie := &http.Cookie{
				Name:     "session",
				Value:    user_uuid.String(),
				HttpOnly: true,
				Path:     "/",
				Expires:  expire,
				MaxAge:   86400,
			}

			session_uid_cookie := &http.Cookie{
				Name:     "session_uid",
				Value:    strconv.Itoa(retUser.Id),
				HttpOnly: true,
				Path:     "/",
				Expires:  expire,
				MaxAge:   86400,
			}

			session_username_cookie := &http.Cookie{
				Name:     "session_username",
				Value:    retUser.Username,
				HttpOnly: true,
				Path:     "/",
				Expires:  expire,
				MaxAge:   86400,
			}

			http.SetCookie(res, session_cookie)
			http.SetCookie(res, session_uid_cookie)
			http.SetCookie(res, session_username_cookie)

			_, err = ctrl.Db.Exec("insert into sessions (user_id, token, created_at) values ($1, $2, $3)", &retUser.Id, user_uuid.String(), current_time)

			if err != nil {
				log.Println(err)
			}

			log.Println("\nUser: %s, has logged in with Session ID UUID: '%s'", retUser.Name, user_uuid)
		} else {
			session_cookie, err := req.Cookie("session")

			if err == nil {
				log.Println("\nUser: %s, is ALREADY logged in with Session ID UUID: '%s'", retUser.Name, session_cookie.Value)
			}
		}

	} else {
		// Notify the user that the password && username do not match
	}

	http.Redirect(res, req, "/", 302)
}

func (ctrl *SessionController) Delete(res http.ResponseWriter, req *http.Request) {
	session_cookie, err := req.Cookie("session")

	if err == nil {
		_, err := ctrl.Db.Exec("delete from sessions where token = $1", session_cookie.Value)

		if err != nil {
			log.Println(err)
		}

		session_cookie = &http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -10,
			Expires:  time.Now(),
		}

		session_uid_cookie := &http.Cookie{
			Name:     "session_uid",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -10,
			Expires:  time.Now(),
		}

		session_username_cookie := &http.Cookie{
			Name:     "session_username",
			Value:    "",
			HttpOnly: true,
			MaxAge:   -10,
			Expires:  time.Now(),
		}

		http.SetCookie(res, session_cookie)
		http.SetCookie(res, session_uid_cookie)
		http.SetCookie(res, session_username_cookie)
	}

	http.Redirect(res, req, "/", 302)
}
