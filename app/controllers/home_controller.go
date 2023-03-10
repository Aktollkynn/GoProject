package controllers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

func register(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "register")
}

func save_user(w http.ResponseWriter, r *http.Request) {
	Fname := r.FormValue("first_name")
	Lname := r.FormValue("last_name")
	Email := r.FormValue("email")
	Password := r.FormValue("password")

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/shop?sslmode=disable")
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
	http.Redirect(w, r, "/home_page", http.StatusSeeOther)
}
func Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "login")
}
func Home_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "home_page")
}
func HandlerRequest() {
	http.HandleFunc("/home_page/", Home_page)
	http.HandleFunc("/login/", Login)
	http.HandleFunc("/save_user/", save_user)
	http.HandleFunc("/register/", register)
	http.ListenAndServe(":9000", nil)
}
