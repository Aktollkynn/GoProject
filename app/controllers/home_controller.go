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

// ------------RegisterAuth------------------------
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

	insert, err := db.Query(fmt.Sprintf("INSERT INTO users (first_name, last_name, email, password) VALUES('%s', '%s', '%s', '%s')", Fname, Lname, Email, hashedPassword))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer insert.Close()

	fmt.Println("Successfully Registered")
	http.Redirect(w, r, "/login?alert=success1", http.StatusSeeOther)
	fmt.Fprintf(w, "<script>alert('You are  sucsesfully!')</script>")

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
	ID        int
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

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Rating      float64
	Comments    []Commenting
	Quantity    int
}

// func productDetailHandler(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")

// 	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	var p Product
// 	err = db.QueryRow("SELECT id, name, description, price, COALESCE((SELECT AVG(rating) FROM ratings WHERE product_id = $1), 0) as avg_rating FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Rating)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	rows, err := db.Query("SELECT id, product_id, comment FROM commenting WHERE product_id = $1", id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var c Commenting
// 		err := rows.Scan(&c.ID, &c.ProductID, &c.Comment)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		p.Comments = append(p.Comments, c)
// 	}

//		tmpl := template.Must(template.ParseFiles("templates/product_detail.html"))
//		tmpl.ExecuteTemplate(w, "product_detail", p)
//	}
func getUserIDFromSession(r *http.Request) (int, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return 0, err
	}
	user, ok := session.Values["user"].(User)
	if !ok {
		return 0, errors.New("User information not found in session")
	}

	return user.ID, nil
}

type Rating struct {
	ID        int
	ProductID int
	Rating    float64
}

func rateProductHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromSession(r)
	if err != nil {
		log.Printf("Error getting user ID from session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		log.Printf("Error converting product ID to integer: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	rating, err := strconv.ParseFloat(r.FormValue("rating"), 64)
	if err != nil {
		log.Printf("Error parsing rating value: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO ratings (user_id, product_id, rating, created_at) VALUES ($1, $2, $3,now())", userID, productID, rating)
	if err != nil {
		log.Printf("Error inserting rate into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product_detail?id="+strconv.Itoa(productID), http.StatusSeeOther)
}

// --------Comment---------------------------
type Commenting struct {
	ID        int
	Comment   string
	ProductID int
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM commenting WHERE product_id = $1", productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	comments := []Commenting{}
	for rows.Next() {
		var c Commenting
		err := rows.Scan(&c.ID, &c.Comment, &c.ProductID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}
	t, err := template.ParseFiles("templates/product_detail.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, comments)
}
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	comment := r.FormValue("comment")
	productIDStr := r.FormValue("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO commenting (comment, product_id) VALUES ($1, $2)", comment, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/product_detail/?id=%d", productID), http.StatusFound)
}

// --------Publishing---------------------------
func AddProductHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/add_product.html"))
	tmpl.Execute(w, nil)
}
func InsertProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		description := r.FormValue("description")
		price := r.FormValue("price")
		db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		_, err = db.Exec("INSERT INTO products (name, description, price) VALUES ($1, $2, $3)", name, description, price)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/add_product/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/add_product/", http.StatusSeeOther)
}

// --------DataInfo-------------------
type DataInfo struct {
	Users      []User
	Products   []Product
	Commenting []Commenting
}

func Data_info(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("postgres", "postgres://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT first_name, last_name, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err = db.Query("SELECT id, name, description, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err = db.Query("SELECT id, comment, product_id FROM commenting")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commenting []Commenting
	for rows.Next() {
		var comment Commenting
		if err := rows.Scan(&comment.ID, &comment.Comment, &comment.ProductID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commenting = append(commenting, comment)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := DataInfo{
		Users:      users,
		Products:   products,
		Commenting: commenting,
	}
	t, err := template.ParseFiles("templates/data_info.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Welcome running")

	t, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.Execute(w, "welcome")
}
func BuyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cardNumber := r.FormValue("card_number")
	var productName string
	err = db.QueryRow("SELECT name FROM products WHERE id = $1", productID).Scan(&productName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the bought product into the "orders" table
	_, err = db.Exec("INSERT INTO orders (product_id, product_name, quantity, card_number) VALUES ($1, $2, $3, $4)", productID, productName, quantity, cardNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	http.Redirect(w, r, "/success?alert=success1", http.StatusSeeOther)
	fmt.Fprintf(w, "<script>alert('You are  sucsesfully!')</script>")
}

func SuccessHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/success.html"))
	tmpl.ExecuteTemplate(w, "success", nil)

}
func productDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	db, err := sql.Open("postgres", "postgresql://postgres:justice@localhost:5432/shop?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var p Product
	err = db.QueryRow("SELECT id, name, description, price, COALESCE((SELECT AVG(rating) FROM ratings WHERE product_id = $1), 0) as avg_rating FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Rating)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, product_id, comment FROM commenting WHERE product_id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Commenting
		err := rows.Scan(&c.ID, &c.ProductID, &c.Comment)
		if err != nil {
			log.Fatal(err)
		}
		p.Comments = append(p.Comments, c)
	}

	tmpl := template.Must(template.ParseFiles("templates/product_detail.html"))
	tmpl.ExecuteTemplate(w, "product_detail", p)
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
	http.HandleFunc("/add_product/", AddProductHandler)
	http.HandleFunc("/product_detail/", productDetailHandler)
	http.HandleFunc("/insert_product/", InsertProductHandler)
	http.HandleFunc("/rate_product/", rateProductHandler)
	http.HandleFunc("/profile/", Profile)
	http.HandleFunc("/edit_profile/", EditProfile)
	http.HandleFunc("/update_profile/", UpdateProfile)
	http.HandleFunc("/data_info/", Data_info)
	http.HandleFunc("/welcome/", Welcome)
	http.HandleFunc("/buy/", BuyHandler)
	http.HandleFunc("/success/", SuccessHandler)

	http.ListenAndServe("localhost:8000", nil)

}
