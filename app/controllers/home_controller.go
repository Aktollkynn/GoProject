package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("name")

	db, err := sql.Open("postgres", "postgresql://postgres:online@localhost:5432/shop?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, description, price FROM products WHERE name LIKE $1", "%"+query+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var results []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
			log.Fatal(err)
		}
		results = append(results, p)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles("templates/search_results.html"))
	tmpl.ExecuteTemplate(w, "search_results", results)
}

var store = sessions.NewCookieStore([]byte("super-secret-key"))

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

	db, err := sql.Open("postgres", "postgresql://postgres:online@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (id, first_name, last_name, email, password) VALUES(DEFAULT, '%s', '%s', '%s', '%s')", Fname, Lname, Email, Password))

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

// New 00:50
var (
	EmailNew    string
	PasswordNew string
)

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("***LoginAuthHandler running***")
	r.ParseForm()
	Email := r.FormValue("email")
	Password := r.FormValue("password")
	fmt.Println(":: Your input ->>  Email:", Email, "Password:", Password)
	EmailNew = Email
	PasswordNew = Password
	db, err := sql.Open("postgres", "postgresql://postgres:online@localhost:5432/shop?sslmode=disable")
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
func Home_page(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("postgres", "postgresql://postgres:online@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, description, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
func HandlerRequest() {
	http.HandleFunc("/home_page/", Home_page)
	http.HandleFunc("/login/", Login)
	http.HandleFunc("/loginauth/", LoginAuth)
	http.HandleFunc("/registerauth/", RegisterAuth)
	http.HandleFunc("/register/", Register)
	http.HandleFunc("/slogout/", SessionLogout)
	http.HandleFunc("/search/", searchHandler)
	http.ListenAndServe("localhost:8000", nil)
}
