package controllers

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gowebapp/helpers"
	"log"
	"net/http"
)

type GeneralController struct {
	Db *sql.DB
}

func NewGeneralController(DBCon *sql.DB) *GeneralController {
	return &GeneralController{DBCon}
}

func (ctrl *GeneralController) Favicon(res http.ResponseWriter, req *http.Request) {
	log.Printf("\nGetting the '%s'.\n", req.URL.Path)
}

func (ctrl *GeneralController) Home(res http.ResponseWriter, req *http.Request) {
	pageData := map[string]interface{}{
		"Title":      "Home",
		"BodyHeader": "Welcome to the Starting Block",
		"Paragraph":  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.",
	}

	user_id, err := req.Cookie("session_uid")

	if err == nil {
		pageData["isUserLoggedIn"] = true

		if foundTweets, userTweets := helpers.GetTweets(user_id.Value, ctrl.Db); foundTweets == true {
			pageData["foundTweets"] = true
			pageData["Tweets"] = userTweets
		} else {
			pageData["foundTweets"] = false
		}
	} else {
		log.Println(err)
		pageData["isUserLoggedIn"] = false
	}

	tpl.ExecuteTemplate(res, "index.html", pageData)
}
