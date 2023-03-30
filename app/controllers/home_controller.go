package controllers

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
}

// --------------------------------------------------
func searchHandler(w http.ResponseWriter, r *http.Request) {
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	query := r.URL.Query().Get("name")

	// Add filter conditions
	var filterConditions []string
	if minPrice != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("price >= %s", minPrice))
	}
	if maxPrice != "" {
		filterConditions = append(filterConditions, fmt.Sprintf("price <= %s", maxPrice))
	}
	filterClause := ""
	if len(filterConditions) > 0 {
		filterClause = "AND " + strings.Join(filterConditions, " AND ")
	}
	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, name, description, price FROM products WHERE name LIKE $1 "+filterClause, "%"+query+"%")

	// rows, err := db.Query("SELECT id, name, description, price FROM products WHERE name LIKE $1", "%"+query+"%")
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

func init() {
	gob.Register(User{})
}

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register running")

	t, err := template.ParseFiles("templates/register.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "register")
}

// --------------------------------------------------
func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterAuthHandler running")

	Fname := r.FormValue("first_name")
	Lname := r.FormValue("last_name")
	email := r.FormValue("email")
	Password := r.FormValue("password")

	if !validateemail(email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if err := ValidateRegistrationForm(Fname, Lname, email, Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (id, first_name, last_name, email, password) VALUES(DEFAULT, '%s', '%s', '%s', '%s')", Fname, Lname, email, hashedPassword))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer insert.Close()

	fmt.Println("Successfully Registered")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// --------------------------------------------------
func ValidateRegistrationForm(fname, lname, email, password string) error {

	if fname == "" || lname == "" {
		return errors.New("Firstname and last name are required")
	}

	if password == "" {
		return errors.New("Password is required")
	}
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}

	return nil
}

// --------------------------------------------------
func validateemail(email string) bool {
	regex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]{2,}$`)
	return regex.MatchString(email)
}

type User struct {
	FirstName string
	LastName  string
	email     string
}

// --------------------------------------------------
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login running")
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "login")
}

// --------------------------------------------------
func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// --------------------------------------------------
func LoginAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginAuthHandler running")
	r.ParseForm()

	email := r.FormValue("email")
	Password := r.FormValue("password")

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var fname, lname, hashedPassword string
	err = db.QueryRow("SELECT first_name, last_name, password FROM users WHERE email = $1", email).Scan(&fname, &lname, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintln(w, "Password or email incorrect!")
			return
		}
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	if err != nil {
		fmt.Fprintln(w, "Password or email incorrect!")
		return
	}

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{
		FirstName: fname,
		LastName:  lname,
		email:     email,
	}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home_page", http.StatusSeeOther)
}

// --------------------------------------------------
// var products = []Product{}
type HomePageData struct {
	User     User
	Products []Product
}

func Home_page(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, ok := session.Values["user"].(User)
	if !ok {
		// session doesn't exist or user information is not stored
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
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

	t, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"User":     user,
		"Products": products,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandlerRequest() {

	http.HandleFunc("/home_page/", Home_page)
	http.HandleFunc("/login/", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/loginauth/", LoginAuth)
	http.HandleFunc("/registerauth/", RegisterAuth)
	http.HandleFunc("/register/", Register)
	http.HandleFunc("/search/", searchHandler)
	http.ListenAndServe("localhost:8000", nil)
}
