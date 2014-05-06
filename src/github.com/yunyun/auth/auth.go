package auth

import (
  "fmt"
  //"time"
  "net/http"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/sessions"
  "github.com/yunyun/db"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "code.google.com/p/go.crypto/bcrypt"
)
const (
  PAGE_LOGIN = "/login/"
  PAGE_ERROR = "/error/"
  FUNC_LOGIN = "/login"
  FUNC_LOGOUT = "/logout"    
  FUNC_REGISTER = "/register"
  FUNC_HOME = "/"
)
/* 
 * user functions
 */
func clear(b []byte) {
    for i := 0; i < len(b); i++ {
        b[i] = 0;
    }
}
func Crypt(password []byte) ([]byte, error) {
    defer clear(password)
    return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
func AddUser(db_ *mgo.Database, user string, email string, pass string) (error) {
  // insert user
  pwHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
  db.IfPanic(err)
  return db.InsertUser(user, email, string(pwHash), db_)
}
func LoginUser(db_ *mgo.Database, name string, pass string) (bson.ObjectId, error)  {
  // query user
  user, err := db.GetUser(name, db_)
  if err != nil {
    fmt.Println("error getting user " + name)
    return "", err
  }
  return user.Id, bcrypt.CompareHashAndPassword([]byte(user.Pw_hash), []byte(pass))
}
/*
 * Web functions
 */
func RequireLogin(rw http.ResponseWriter, req *http.Request, s sessions.Session,  
                  db_ *mgo.Database, c martini.Context) {
  user, err := db.GetUserById(bson.ObjectIdHex(s.Get("userId").(string)), db_)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, req, PAGE_LOGIN, http.StatusFound)
    return
  }
  c.Map(user)
}
func Register(rw http.ResponseWriter, r *http.Request, db_ *mgo.Database) {
  name, email, password := r.FormValue("name"), r.FormValue("email"), r.FormValue("password")
  
  err := AddUser(db_, name, email, password)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, r, PAGE_ERROR, http.StatusFound)
  }
  http.Redirect(rw, r, PAGE_LOGIN, http.StatusFound)
}
func Logout(rw http.ResponseWriter, req *http.Request, s sessions.Session) {
  s.Delete("userId")
  http.Redirect(rw, req, PAGE_LOGIN, http.StatusFound)
}
func PostLogin(rw http.ResponseWriter, req *http.Request, db_ *mgo.Database, s sessions.Session) {
  user, pass := req.FormValue("email"), req.FormValue("password")
  id, err := LoginUser(db_, user, pass)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, req, PAGE_LOGIN, http.StatusFound)
  }
  s.Set("userId", id.Hex())
  http.Redirect(rw, req, FUNC_HOME, http.StatusFound)
}