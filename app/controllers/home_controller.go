package controllers

import (
	"html/template"
	"net/http"

	// "gopkg.in/kataras/go-serializer.v0/data"  05.03.23
	
)

//	func Home(w http.ResponseWriter, r *http.Request) {
//		tmpl, err := template.ParseFiles("templates/home.tmpl")
//		if err != nil {
//			// handle error
//		}
//
//		_ = render.HTML(w, http.StatusOK, "home", map[string]interface{}{
//			"title": "Home Title",
//			"body":  "Home Description",
//		})
//
// }

type Register struct {
	First_name       string
	Last_name        string
	Email            string
	Password         string
	Confirm_password string
}

var (
	tem = template.Must(template.ParseFiles("templates/register.html"))
)

func Registerr(w http.ResponseWriter, r *http.Request) {
	data := Register{
		First_name:       r.FormValue("first_name"),
		Last_name:        r.FormValue("last_name"),
		Email:            r.FormValue("email"),
		Password:         r.FormValue("password"),
		Confirm_password: r.FormValue("confirm_password"),
	}
	tem.Execute(w, data)
}
func sign_in_page(w http.ResponseWriter, r *http.Request) {
    const name = "aaaaa"
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(w, name)
}

func Home_page(w http.ResponseWriter, r *http.Request) {
	const name = "aaaa"
	t, _ := template.ParseFiles("templates/home_page.html")
	t.Execute(w, name)
}

func HandlerRequest() {
	http.HandleFunc("/home_page/", Home_page)
    http.HandleFunc("/sign_in/",  sign_in_page)
	http.HandleFunc("/register/", Registerr)
	http.ListenAndServe(":9000", nil)
}
