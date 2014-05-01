// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"
  "database/sql"
  "net/http"
  "github.com/yunyun/auth"
  "github.com/yunyun/kotoba"
  //"github.com/yunyun/news"  
  "github.com/yunyun/db"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"  
  "github.com/martini-contrib/sessions"
  "strconv"
  //"time"
)
var pageLogin string = "/login/"
var pageError string = "/error/"
var pageAdd string = "/addWord"
var pageReviews string = "/doReviews"
var funcLogin string = "/login"
var funcLogout string = "/logout"
var funcRegister string = "/register"
var funcKotoba string = "/kotoba/:id"
var funcHome string = "/"
var funcAdd string = "/add"
var funcCheck string = "/reviewed/:id"
var apiReviews string = "/apiReviews/:apikey"

type User struct {
  Id int
  Name string
  Email string
}

func SetupDB() *sql.DB {
  datab, err := db.GetDB()
  if err != nil {
    panic(err)
  }
  return datab
}

/*
 * Main function
 */
func main() {
  m := martini.Classic()
  
  //sessions
  store := sessions.NewCookieStore([]byte("yyt726ddd318"))
  m.Use(sessions.Sessions("yunyun", store))
  
  m.Map(SetupDB())
  //render html templates from directory
  m.Use(render.Renderer(render.Options{
    Layout: "layout",
  }))
  
  m.Post(funcLogin, PostLogin)  
  m.Post(funcRegister, Register)
  m.Post(funcAdd, RequireLogin, AddKotoba)
  m.Post(funcCheck, RequireLogin, CheckReview)
  m.Get(funcLogout, Logout)
  m.Get(funcHome, RequireLogin, Home)  
  m.Get(funcKotoba, RequireLogin, Kotoba)
  m.Get(pageAdd, RequireLogin, AddWord)
  m.Get(pageReviews, RequireLogin, DoReviews)
  m.Get(apiReviews, RequireLogin, GetReviews)
 
	m.Run()
}

func RequireLogin(rw http.ResponseWriter, req *http.Request, s sessions.Session,  
                  db *sql.DB, c martini.Context) {
  
  user := &User{}
  id, name, email, err := auth.GetUser(db, s.Get("userId"))
  user.Name = name
  user.Email = email
  user.Id = id
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, req, pageLogin, http.StatusFound)
    return
  }
  c.Map(user)
}

func Logout(rw http.ResponseWriter, req *http.Request, s sessions.Session) {
  s.Delete("userId")
  http.Redirect(rw, req, pageLogin, http.StatusFound)
}
func Register(rw http.ResponseWriter, r *http.Request, db *sql.DB) {
  name, email, password := r.FormValue("name"), r.FormValue("email"), r.FormValue("password")
  
  err := auth.AddUser(db, name, email, password)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, r, pageError, http.StatusFound)
  }
  http.Redirect(rw, r, pageLogin, http.StatusFound)
}
func AddKotoba(rw http.ResponseWriter, r *http.Request, u *User, db *sql.DB) {
  word, hatsuon, hatsuon_, imi, imi_ := r.FormValue("word"), r.FormValue("hatsuon"), r.FormValue("hatsuon_"), r.FormValue("imi"), r.FormValue("imi_")

  k := kotoba.NewDefaultLevelKotoba(u.Id, word, hatsuon, hatsuon_, imi, imi_)
  err := k.Save(db)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, r, pageError, http.StatusFound)
  }
  http.Redirect(rw, r, funcHome, http.StatusFound)
}
func CheckReview(params martini.Params, rw http.ResponseWriter, r *http.Request, u *User, db *sql.DB) {
  checked := r.FormValue("checked")
  id, _ := strconv.Atoi(params["id"])
  k := kotoba.GetKotoba(id, db)
  fmt.Println(k)
  if checked == "true" {
    k.IncLevel()
    k.Update(db)
  } else {
    k.DecLevel()
    k.Update(db)
  }
  fmt.Println(k)
  http.Redirect(rw, r, pageReviews, http.StatusFound)
}

func GetReviews(params martini.Params, r render.Render, u *User, db *sql.DB) {
  r.JSON(200, map[string]interface{}{"hello" : params["apikey"]})
}
func PostLogin(rw http.ResponseWriter, req *http.Request, db *sql.DB, s sessions.Session) {
  user, pass := req.FormValue("email"), req.FormValue("password")
  id, err := auth.LoginUser(db, user, pass)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, req, pageLogin, http.StatusFound)
  }
  s.Set("userId", id)
  http.Redirect(rw, req, funcHome, http.StatusFound)
}

func AddWord(r render.Render, u *User) {
  r.HTML(200, "add", nil)
}

func Home(rw http.ResponseWriter, req *http.Request, r render.Render, u *User, db *sql.DB) {
 
  k, err := kotoba.GetAllKotoba(u.Id, db)
  if err != nil {
    fmt.Println(err)
    http.Redirect(rw, req, pageError, http.StatusFound)
  }
  r.HTML(200, "home", k)
}

func DoReviews(r render.Render, u *User, db *sql.DB) {
  k := kotoba.GetFirstReviewKotoba(u.Id, db)
  r.HTML(200, "review", k)
}

func Kotoba(params martini.Params, r render.Render, u *User, db *sql.DB) {
  id, _ := strconv.Atoi(params["id"])
  k := kotoba.GetKotoba(id, db)
  r.HTML(200, "kotoba", k)
}
