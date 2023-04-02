// `home_controller.go`
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

// -----------------Search------------------
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

// ------------RegisterAuth-------------------------
func RegisterAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterAuthHandler running")

	Fname := r.FormValue("first_name")
	Lname := r.FormValue("last_name")
	Email := r.FormValue("email")
	Password := r.FormValue("password")

	if !validateemail(Email) {
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

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
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
	http.Redirect(w, r, "/login?alert=success1", http.StatusSeeOther)
	fmt.Fprintf(w, "<script>alert('You are registered sucsesfully!')</script>")

}

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

func validateemail(Email string) bool {
	regex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]{2,}$`)
	return regex.MatchString(Email)
}

type User struct {
	FirstName string
	LastName  string
	Email     string
}

// ---------------Login----------------------------------
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login running")
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "login")
}

// -------------Logout-------------------------
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

// --------------Loginaut----------------------------
func LoginAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginAuthHandler running")
	r.ParseForm()

	Email := r.FormValue("email")
	Password := r.FormValue("password")

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var fname, lname, hashedPassword string
	err = db.QueryRow("SELECT first_name, last_name, password FROM users WHERE email = $1", Email).Scan(&fname, &lname, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Email address is not registered", http.StatusUnauthorized)
			return
		}
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(Password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
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
		Email:     Email,
	}
	session.Values["user"] = user
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/home_page", http.StatusSeeOther)
}

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

// --------Profile---------------------------

type UserInfo struct {
	FirstName string
	LastName  string
	Email     string
}

func Profile(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := UserInfo{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	err = t.Execute(w, userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func EditProfile(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseFiles("templates/edit_profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := UserInfo{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	err = t.Execute(w, userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func ValidateEditProfileForm(fname, lname, email, password string) error {

	if fname == "" || lname == "" {
		return errors.New("Firstname and last name are required")
	}

	if email == "" {
		return errors.New("Email is required")
	}

	return nil
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	Fname := r.FormValue("first_name")
	Lname := r.FormValue("last_name")
	Email := r.FormValue("email")
	Password := r.FormValue("password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)

	if !validateemail(Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if err := ValidateEditProfileForm(Fname, Lname, Email, ""); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Query(fmt.Sprintf("UPDATE users SET first_name='%s', last_name='%s', email='%s', password='%s' WHERE email='%s'", Fname, Lname, Email, hashedPassword, user.Email))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{
		FirstName: Fname,
		LastName:  Lname,
		Email:     Email,
	}
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile?alert=success", http.StatusSeeOther)
	fmt.Fprintf(w, "<script>alert(' Your information has been changed successfully!.')</script>")
}

// --------HandlerRequest-------------------

func HandlerRequest() {

	http.HandleFunc("/home_page/", Home_page)
	http.HandleFunc("/login/", Login)
	http.HandleFunc("/logout/", Logout)
	http.HandleFunc("/loginauth/", LoginAuth)
	http.HandleFunc("/registerauth/", RegisterAuth)
	http.HandleFunc("/register/", Register)
	http.HandleFunc("/search/", searchHandler)

	http.HandleFunc("/profile/", Profile)
	http.HandleFunc("/edit_profile/", EditProfile)
	http.HandleFunc("/update_profile/", UpdateProfile)

	http.ListenAndServe("localhost:8000", nil)

}
