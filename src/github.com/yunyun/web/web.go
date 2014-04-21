package util

import (
  "fmt"
	"html/template"
  "net/http"
	"regexp"
  "time"
  "github.com/dgrijalva/jwt-go"
)
/*
 * Type definitions
 */
type UserData struct {
  User string
}

/*
 * Constants definitions
 */
const (
  JWT_SIGNING_KEY = "!di%06Z&d#;ldPW@d"
)

/*
 * define variables
 */
var templates = template.Must(template.ParseFiles("web/login.html", "web/account.html", "web/register.html"))
var validPath = regexp.MustCompile("^/(login|account|register)/?([a-zA-Z0-9]*)$")

/* 
 * webservice functions
 */
func renderTemplate(w http.ResponseWriter, tmpl string, data UserData) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
//func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		m := validPath.FindStringSubmatch(r.URL.Path)
//		if m == nil {
//			http.NotFound(w, r)
//			return
//		}
//		fn(w, r, m[2])
//	}
//}
func createNewJWT(user string) string {
  token := jwt.New(jwt.GetSigningMethod("HS256"))
  token.Claims["user"] = user
  token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
  tokenString, err := token.SignedString([]byte(JWT_SIGNING_KEY))
  if err != nil {
    panic(err.Error())
  }
  return tokenString
}
func validateJWT(jwtKey string) {
//  token, err := jwt.Parse(params["token"], func(token *jwt.Token) ([]byte, error) {
//    return []byte(SecretKey), nil
//  })
//  if err == nil && token.Valid {
//    return "User id: " + token.Claims["userid"].(string)
//  } else {
//    return "Invalid"
//  }
}
func AccountHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("account page")
  user := r.FormValue("user")
  pass := r.FormValue("pass")
  //fmt.Println("user: '" + user + "'")
  //fmt.Println("pass: '" + pass + "'")
  
  err := loginUser(user, pass)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	  http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
  token := createNewJWT(user)
  fmt.Println(token)
	renderTemplate(w, "account", UserData{User: user})
}
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("redirect page")
	  http.Redirect(w, r, "/login", http.StatusFound)
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("register page")
	renderTemplate(w, "register", UserData{User: ""})
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("login page")
	renderTemplate(w, "login", UserData{User: ""})
}