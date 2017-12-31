
package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
	"github.com/satori/go.uuid"
)


type User struct {
    Name 	string 		  `bson:"name,omitempty"`
    Email 	string 		  `bson:"email,omitempty"`
    Hash 	string 		  `bson:"hash,omitempty"`
}

type Session struct {
	Email 	string 	`bson:"email,omitempty"`
	Uuid 	string 	`bson:"uuid,omitempty"`
}


var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}


func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

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

    fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}




func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm();
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	retUser := User{}

	session, err := mgo.Dial("mongodb://127.0.0.1:27017")	
	if err != nil {
		panic(err)
	}
	defer session.Close()

	col := session.DB("webdev").C("users")
	err = col.Find(bson.M{"email": email}).One(&retUser)

	if err != nil {
		panic(err)
	}

	// fmt.Println("\nResult: ", retUser)

	pwd_match := bcrypt.CompareHashAndPassword([]byte(retUser.Hash), []byte(password))

	if pwd_match == nil {
		session_cookie, err := r.Cookie("session")

		if err != nil {
			user_uuid := uuid.NewV4()
			session_cookie = &http.Cookie{
				Name: "session",
				Value: user_uuid.String(),
				HttpOnly: true,
			}
			http.SetCookie(w, session_cookie)

			sCol := session.DB("webdev").C("sessions")
			sCol.Insert(Session{email, user_uuid.String()})

			fmt.Printf("\nUser: %s, has logged in with Session ID UUID: '%s'", retUser.Name, user_uuid)	
		} else {
			session_cookie, err := r.Cookie("session")
			if err == nil {
				fmt.Printf("\nUser: %s, is ALREADY logged in with Session ID UUID: '%s'", retUser.Name, session_cookie.Value)
			}
		}

	} else {
		fmt.Printf("\nNAH, YOU TRIED IT. YAH FAILED. GO AGAIN.")
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}


func logout(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Connect to the "sessions" collection in "webdev"
	sCol := session.DB("webdev").C("sessions")

	session_cookie, err := r.Cookie("session")

	if err == nil {
		// Create a temporary variable to hold a "Session" struct for later use
		// tempUUID := Session{}
		// Find ONE document where the "uuid" is the same as the current users session uuid value
		// sCol.Find(bson.M{"uuid": session_cookie.Value}).One(&tempUUID)
		// fmt.Println("\nFound the session object im about to delete: ", tempUUID)

		// Remove the document where the "uuid" field is the same as the users session UUID value
		sCol.Remove(bson.M{"uuid": session_cookie.Value })

		session_cookie = &http.Cookie{
			Name: "session",
			Value: "",
			HttpOnly: true,
		}
		http.SetCookie(w, session_cookie)
	}

	fmt.Printf("\nRedirecting to the '/' path\n")
	http.Redirect(w, r, "/", 303)
}










