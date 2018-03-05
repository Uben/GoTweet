package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template
var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=devtest dbname=gowebapp password=password sslmode=disable")

	if err != nil {
		panic(err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	gmux := mux.NewRouter()

	gmux.HandleFunc("/register", user_register).Methods("GET")
	gmux.HandleFunc("/register", register).Methods("POST")

	gmux.HandleFunc("/login", user_login).Methods("GET")
	gmux.HandleFunc("/login", login).Methods("POST")
	gmux.HandleFunc("/logout", isAuth(logout)).Methods("GET")

	gmux.HandleFunc("/settings", update_user).Methods("GET")
	gmux.HandleFunc("/update-user-info", change_user_info).Methods("POST")
	gmux.HandleFunc("/update-user-password", change_user_password).Methods("POST")

	gmux.HandleFunc("/create-tweet", tweet_create).Methods("POST")
	gmux.HandleFunc("/delete-tweet/{tweet_id}", tweet_delete).Methods("POST")

	gmux.HandleFunc("/favicon.ico", handlerIcon).Methods("GET")
	gmux.HandleFunc("/", home).Methods("GET")

	fmt.Printf("About to listen on port :3000. Go to https://127.0.0.1:3000/ (localhost)\n")
	log.Fatal(http.ListenAndServe(":3000", gmux))
}

func handlerIcon(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)
}

func home(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Create map to pass data to template
	pageData := map[string]string{
		"Title":      "Bernin Uben | Base Golang Web App",
		"BodyHeader": "Welcome to the Starting Block",
		"Paragraph":  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.",
	}

	if is_user_logged_in(req) {
		pageData["isUserLoggedIn"] = "true"
	} else {
		pageData["isUserLoggedIn"] = "false"
	}

	// Check if the path is exactly "/" else its a 404 error
	if req.URL.Path != "/" {
		tpl.ExecuteTemplate(res, "404.html", pageData)

		// Else Execute the index template with the 'pageData' data
	} else {
		tpl.ExecuteTemplate(res, "index.html", pageData)
	}
}
