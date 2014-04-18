// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
  "code.google.com/p/go.crypto/bcrypt"
)
/*
 * Constants definitions
 */
const (
    DB_HOST = "tcp(127.0.0.1:3306)"
    DB_NAME = "yunyun"
    DB_USER = "yunyun"
    DB_PASS = "nuynuy"
)
/*
 * Type definitions
 */
type Page struct {
	User string
	Body  []byte
}
type User struct {
  user string
  pwHash string
}
type DateTime struct {
  year int
  month int
  day int
  hour int
  min int
  sec int
}
type Kotoba struct {
  kotoba string
  imi string
  level int
  review DateTime  
}

/*
 * define variables
 */
var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)
var templates = template.Must(template.ParseFiles("login.html", "account.html", "register.html"))
var validPath = regexp.MustCompile("^/(login|account|register)/?([a-zA-Z0-9]*)$")

/* 
 * test functions
 */
func (p *Page) save() error {
	filename := p.User + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{User: title, Body: body}, nil
}

/* 
 * database functions
 */
func printRows(rows *sql.Rows) {
  cols, _ := rows.Columns()
  n := len(cols)
  
  for i := 0; i < n; i++ {
    fmt.Print(cols[i], "\t")
  }
  fmt.Println()
  
  var fields []interface{}
  for i := 0; i < n; i++ {
    fields = append(fields, new(string))
  }
  for rows.Next() {
    rows.Scan(fields...)
    for i := 0; i < n; i++ {
      fmt.Print(*(fields[i].(*string)), "\t")
    }
    fmt.Println()
  }
  fmt.Println()
}
func printTable(table string) {
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
  rows, err := db.Query("SELECT * FROM " + table);
  if err != nil {
    panic(err.Error()) // proper error handling instead of panic in your app
  }
  printRows(rows)
}
 
/* 
 * webservice functions
 */
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
func accountHandler(w http.ResponseWriter, r *http.Request, title string) {
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
	renderTemplate(w, "account", &Page{User: user, Body: []byte("")})
}
func registerHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{User: title}
	}
	renderTemplate(w, "register", p)
}
func loginHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{User: title}
	}
	renderTemplate(w, "login", p)
}

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

/* 
 * kotoba functions
 */


/*
 * Main function
 */
func main() {
  //printTable("users")
  //printTable("kotoba")
  //err := addUser(user, pass)
  //err := loginUser(user, pass)
  //if err != nil {
  //  fmt.Println(err.Error());
  //}

  flag.Parse()
	http.HandleFunc("/account/", makeHandler(accountHandler))
	http.HandleFunc("/login", makeHandler(loginHandler))
	http.HandleFunc("/register", makeHandler(registerHandler))

	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":3000", nil)
}
