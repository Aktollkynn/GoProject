package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)



func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***Register running***")

	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "register")
}

func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***RegisterAuthHandler running***")

	Fname := r.FormValue("first_name")
	Lname := r.FormValue("last_name")
	Email := r.FormValue("email")
	Password := r.FormValue("password")

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (first_name, last_name, email, password) VALUES('%s', '%s', '%s', '%s')", Fname, Lname, Email, Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer insert.Close()
	fmt.Printf("~~~New user altered :[first_mame: '%s', last_name:'%s', email:'%s', password: '%s']", Fname, Lname, Email, Password)

	fmt.Println("==Successfully Registered==")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***Login running***")
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "login")
}

// New 12.02.23
var (
	EmailNew    string
	PasswordNew string
)
var store = sessions.NewCookieStore([]byte("super-secret-key"))

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***LoginAuthHandler running***")
	r.ParseForm()
	Email := r.FormValue("email")
	Password := r.FormValue("password")
	fmt.Println(":: Your input ->>  Email:", Email, "Password:", Password)
	EmailNew = Email
	PasswordNew = Password
	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	//-----
	rows, err := db.Query("SELECT email, password FROM users order by id desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	fmt.Println("---All users from DB:")
	countuser := 0
	for rows.Next() {
		err := rows.Scan(&Email, &Password)
		if EmailNew == Email && PasswordNew == Password {
			EmailNew = Email
			fmt.Println("--> Email and Password Correct! WELCOME ####")
			http.Redirect(w, r, "/home_page", http.StatusSeeOther)
			countuser += 1

		}

		if err != nil {
			panic(err)
		}
		fmt.Println("\n:::", Email, Password)
	}

	EmailNew = ""
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if countuser <= 0 {
		fmt.Println(w, "Password or Email incorrect, Try again!")
		fmt.Fprintf(w, "Password or Email incorrect, Try again!")
		return
	}

}

func SessionLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/login", 302)
}
func SessionLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = true
	session.Values["email"] = r.FormValue("email")
	session.Values["password"] = r.FormValue("password")
	session.Save(r, w)
	http.Redirect(w, r, "/welcome", 302)
}
func Home_page(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***Home running***")

	t, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "home_page")
}
func HandlerRequest() {
	http.HandleFunc("/home_page/", Home_page)
	http.HandleFunc("/login/", Login)
	http.HandleFunc("/loginauth/", LoginAuth)
	http.HandleFunc("/registerauth/", RegisterAuth)
	http.HandleFunc("/register/", Register)
	http.HandleFunc("/slogout/", SessionLogout)
	http.ListenAndServe("localhost:8000", nil)
}
