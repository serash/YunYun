package util

import (
  "log"
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
func addUser(user string, email string, pass string) (error) {
  // get database
  dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
  db, err := sql.Open("mysql", dsn)
  if err != nil {
     panic(err.Error())
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
  }
  // insert user
  pwHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
  if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
  }
  _, err = db.Exec("INSERT INTO users(user, email, pw_hash) VALUES(?, ?, ?)", user, email, pwHash)
  return err
}
func loginUser(user string, pass string) (error)  { 
  // get database
  dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
  db, err := sql.Open("mysql", dsn)
  if err != nil {
     panic(err.Error())
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
  }
  // query user
  var pwHash string
  err = db.QueryRow("SELECT pw_hash FROM users WHERE (user = ? OR email = ?)", user, user).Scan(&pwHash)
  switch {
    case err == sql.ErrNoRows:
      log.Printf("No user found.")
    case err != nil:
            log.Fatal(err)
  }
  return bcrypt.CompareHashAndPassword([]byte(pwHash), []byte(pass))
}