package main

import (
	"database/sql"
	"fmt"
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

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/favicon.ico", handlerIcon)
	http.HandleFunc("/", home)

	fmt.Printf("About to listen on port :3000. Go to https://127.0.0.1:3000/ (localhost)\n")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handlerIcon(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)
}

func home(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)

	// Check if user is logged in
	_, userLoggedInInt := is_user_logged_in(req)

	// Create map to pass data to template
	pageData := map[string]string{
		"Title":          "Bernin Uben | Software Developer",
		"BodyHeader":     "Welcome to the Starting Block",
		"Paragraph":      "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Urna cursus eget nunc scelerisque viverra. Tincidunt nunc pulvinar sapien et ligula ullamcorper. Suspendisse potenti nullam ac tortor vitae.",
		"IsUserLoggedIn": string(userLoggedInInt),
	}

	// Check if the path is exactly "/" else its a 404 error
	if req.URL.Path != "/" {
		tpl.ExecuteTemplate(res, "404.gohtml", pageData)
	} else {
		// Execute the template with the 'pageData' data
		tpl.ExecuteTemplate(res, "index.gohtml", pageData)
	}
}
