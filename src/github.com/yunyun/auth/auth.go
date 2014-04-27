package auth

import (
  "fmt"
  "strconv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "code.google.com/p/go.crypto/bcrypt"
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
func addUser(db *sql.DB, user string, email string, pass string) (error) {
  // insert user
  pwHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
  if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
  }
  _, err = db.Exec("INSERT INTO users(user, email, pw_hash) VALUES(?, ?, ?)", user, email, pwHash)
  return err
}
func LoginUser(db *sql.DB, user string, pass string) (string, error)  {
  // query user
  var pwHash string
  var id string
  err := db.QueryRow("SELECT id, pw_hash FROM users WHERE (user = ? OR email = ?)", user, user).Scan(&id, &pwHash)
  if err != nil {
    fmt.Println("error getting user " + user)
    return "", err
  }
  return id, bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(pass))
}

func GetUser(db *sql.DB, user interface{}) (int, string, string, error) {
  var name string
  var email string
  var id_ string
  err := db.QueryRow("SELECT id, user, email from users where id=?", user).Scan(&id_, &name, &email)
  id, _ := strconv.Atoi(id_)
  return id, name, email, err
}