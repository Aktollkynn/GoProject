package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"regexp"
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
	fmt.Println("Register running")

	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "register")
}

func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterAuthHandler running")

	Fname := r.FormValue("first_name")
	Lname := r.FormValue("last_name")
	Email := r.FormValue("email")
	Password := r.FormValue("password")

	if !validateEmail(Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if err := ValidateRegistrationForm(Fname, Lname, Email, Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
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

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (id, first_name, last_name, email, password) VALUES(DEFAULT, '%s', '%s', '%s', '%s')", Fname, Lname, Email, hashedPassword))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer insert.Close()

	fmt.Println("Successfully Registered")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ValidateRegistrationForm(fname, lname, email, password string) error {

	if fname == "" || lname == "" {
		return errors.New("Firs tname and last name are required")
	}

	if password == "" {
		return errors.New("Password is a required")
	}
	if len(password) < 8 {
		return errors.New("Password must be at least 8 symbols")
	}
	return nil
}

func validateEmail(email string) bool {
	regex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]{2,}$`)
	return regex.MatchString(email)
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login running")
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "login")
}

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginAuthHandler running")
	r.ParseForm()
	Email := r.FormValue("email")
	Password := r.FormValue("password")

	db, err := sql.Open("postgres", "postgresql://postgres:online@localhost:5432/shop?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var hashedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = $1", Email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "Password or Email incorrect!")
			return
		}
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	if err != nil {
		fmt.Fprintf(w, "Password or Email incorrect!")
		return
	}
	fmt.Println("Email and Password Correct!")
	http.Redirect(w, r, "/home_page", http.StatusSeeOther)

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
