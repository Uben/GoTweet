
package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	// Id 		bson.ObjectId `bson:"_id,omitempty"`
    Name 	string 		  `bson:"name,omitempty"`
    Email 	string 		  `bson:"email,omitempty"`
    Hash 	string 		  `bson:"hash,omitempty"`
}


var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}


func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)

	fmt.Printf("About to listen on port :3000. Go to https://127.0.0.1:3000/\n")
	log.Fatal(http.ListenAndServe(":3000", nil))
}



func home(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("\nUser accessed the '%s' url path.\n", req.URL.Path)
	// tpl := template.Must(template.ParseFiles("templates/index.gohtml"))

	pageData := map[string]string {
		"Title" : "Bernin Uben | Software Developer",
		"BodyHeader" : "Welcome to the Starting Block",
		"Paragraph" : "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Urna cursus eget nunc scelerisque viverra. Tincidunt nunc pulvinar sapien et ligula ullamcorper. Suspendisse potenti nullam ac tortor vitae.",
	}

	tpl.ExecuteTemplate(res, "index.gohtml", pageData)
}


func register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nUser accessed the '%s' url path.\n", r.URL.Path)
	r.ParseForm();
	name := r.PostFormValue("name")
    email := r.PostFormValue("email")
    password := r.PostFormValue("password")
    confirm_password := r.PostFormValue("confirm_password")
    new_user := User{}

	if password == confirm_password {
    	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    	
    	if err != nil {
    		log.Println(err)
    	}

    	new_user = User{name, email, string(hashPass)}

    	if bcrypt.CompareHashAndPassword(hashPass, []byte(password)) == nil {
			session, err := mgo.Dial("mongodb://127.0.0.1:27017")
			if err != nil {
				panic(err)
			}
			defer session.Close()

			col := session.DB("webdev").C("users")
			col.Insert(new_user)
    	}
    }

	// fmt.Printf("\nUser accessed the '%s' url path.\n", r.URL.Path)
	// fmt.Printf("User Name: '%s'.\n", name)
	// fmt.Printf("User Email: '%s'.\n", email)
	// fmt.Printf("User Password: '%s'.\n", password)
	// fmt.Printf("User Confirm Password: '%s'.\n", confirm_password)

	http.Redirect(w, r, "/", 303)
}


func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	result := User{}

	session, err := mgo.Dial("mongodb://127.0.0.1:27017")		
	if err != nil {
		panic(err)
	}
	defer session.Close()

	col := session.DB("webdev").C("users")

	err = col.Find(bson.M{"email": email}).One(&result)

	if err != nil {
		panic(err)
	}

	fmt.Println("\nResult: ", result)

	pwd_match := bcrypt.CompareHashAndPassword([]byte(result.Hash), []byte(password))

	if pwd_match == nil {
		fmt.Printf("\nUser: %s, has logged in.", result.Name)
	} else {
		fmt.Printf("\nNAH, YOU TRIED IT. YAH FAILED. GO AGAIN.")
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}











