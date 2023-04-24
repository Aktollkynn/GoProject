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
	"strconv"
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
	Rating      float64
}

// -----------------Search------------------
func searchHandler(w http.ResponseWriter, r *http.Request) {
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	query := r.URL.Query().Get("name")

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
	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
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

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (first_name, last_name, email, password) VALUES('%s', '%s', '%s', '%s')", Fname, Lname, Email, hashedPassword))

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

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
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

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
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

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
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

func productDetailHandler(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var p Product
	//err = db.QueryRow("SELECT id, name, description, price FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Description, &p.Price)
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = db.QueryRow("SELECT id, name, description, price, COALESCE((SELECT AVG(rating) FROM ratings WHERE product_id = $1), 0) as avg_rating FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Rating)
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(template.ParseFiles("templates/product_detail.html"))
	tmpl.ExecuteTemplate(w, "product_detail", p)
}
func rateProductHandler(w http.ResponseWriter, r *http.Request) {

	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		log.Printf("Error converting product_id: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		log.Printf("Error converting rating: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO ratings (product_id, rating) VALUES ($1, $2)", productID, rating)
	if err != nil {
		log.Printf("Error inserting rating into database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/product_detail/?id=%d", productID), http.StatusFound)
}

type Commenting struct {
	ID        int
	Comment   string
	Name      sql.NullString
	ProductID int
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve all comments from the database.
	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	comments := []Commenting{}

	for rows.Next() {
		var c Commenting
		err := rows.Scan(&c.ID, &c.Comment, &c.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	t, err := template.ParseFiles("templates/comment.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, comments)
}
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	name := r.FormValue("name")
	var nameNull sql.NullString
	if name != "" {
		nameNull.String = name
		nameNull.Valid = true
	}

	db, err := sql.Open("postgres", "postgresql://postgres:aktolkyn@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO comments (comment,name) VALUES ($1, $2)", comment, nameNull)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/comment", http.StatusSeeOther)
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
	http.HandleFunc("/comment/", CommentHandler)
	http.HandleFunc("/add_comment/", AddCommentHandler)

	http.HandleFunc("/profile/", Profile)
	http.HandleFunc("/edit_profile/", EditProfile)
	http.HandleFunc("/update_profile/", UpdateProfile)
	http.HandleFunc("/product_detail/", productDetailHandler)
	http.HandleFunc("/rate_product/", rateProductHandler)

	http.ListenAndServe("localhost:8000", nil)

}
