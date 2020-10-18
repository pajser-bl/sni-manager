package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"sni-admin/user"
	"sni-manager/article"
	"strconv"
)

func GetUserType(req *http.Request) int {
	il, err := IsLoggedIn(req)
	if err != nil {
		return -1
	}
	return int(il.Type)
}

func indexPageHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/index.html")
}
func loginPageHandler(res http.ResponseWriter, req *http.Request) {
	http.Redirect(res, req, "https://localhost/admin/login", http.StatusSeeOther)
}

func articlesHandler(res http.ResponseWriter, req *http.Request) {
	if GetUserType(req) != 0 && GetUserType(req) != 1 {
		http.Redirect(res, req, "https://localhost/managment", http.StatusForbidden)
	} else {
		Articles, _ := article.GetAllArticles(db)
		t, _ := template.ParseFiles("./static/articles.html")
		if req.Method == "GET" {
			_ = t.Execute(res, Articles)
		} else if req.Method == "POST" {
			un := template.HTMLEscapeString(req.FormValue("username"))
			fn := template.HTMLEscapeString(req.FormValue("firstname"))
			ln := template.HTMLEscapeString(req.FormValue("lastname"))
			pw := template.HTMLEscapeString(req.FormValue("password"))
			typ, _ := strconv.Atoi(template.HTMLEscapeString(req.FormValue("type")))
			ph, _ := bcrypt.GenerateFromPassword([]byte(pw), 14)
			if un != "" && fn != "" && ln != "" && pw != "" {
				_, err := user.Create(db, &user.User{Username: un, FirstName: fn, LastName: ln, PasswordHash: string(ph), Type: uint8(typ)})
				if err != nil {
					panic(err.Error())
				}
				Articles, _ = article.GetAllArticles(db)
				_ = t.Execute(res, Articles)
			} else {
				http.Redirect(res, req, "https://localhost/manager/articles", http.StatusSeeOther)
			}
		}
	}
}

var red *redis.Client
var db *gorm.DB

func main() {

	db = dbConn()
	red = getRedisClient()
	router := mux.NewRouter()
	router.HandleFunc("/", indexPageHandler).Name("home")
	router.HandleFunc("/manager", indexPageHandler).Name("home")
	router.HandleFunc("/home", indexPageHandler).Name("home")
	router.HandleFunc("/index", indexPageHandler).Name("home")
	router.HandleFunc("/login", loginPageHandler).Name("login")
	router.HandleFunc("/articles", articlesHandler).Name("articles")
	//router.HandleFunc("/articles/{id:[0-9]+}", articleHandler).Name("article")

	log.Println("Listening on :5000")
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		log.Fatal(err)
	}
}
