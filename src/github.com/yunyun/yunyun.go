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
  "github.com/yunyun/news"  
  "github.com/yunyun/db"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"  
  "github.com/martini-contrib/sessions"
  "strconv"
  //"time"
)

type User struct {
  Id int
  Name string
  Email string
}

func SetupDB() *sql.DB {
  db, err := db.GetDB()
  if err != nil {
    panic(err)
  }
  return db
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
  
  m.Post("/login", PostLogin)
  m.Get("/logout", func(s sessions.Session) string {
    s.Delete("userId")
    return "Logged out"
  })
  m.Get("/", Home)  
  m.Get("/kotoba", RequireLogin, Kotoba)
	m.Get("/kotoba/:id", RequireLogin, KotobaId)
 
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
    http.Redirect(rw, req, "/", http.StatusFound)
    return
  }
  c.Map(user)
}

func PostLogin(req *http.Request, db *sql.DB, s sessions.Session) (int, string){
  user, pass := req.FormValue("email"), req.FormValue("password")
  id, err := auth.LoginUser(db, user, pass)
  if err != nil {
    return 401, "Not Authorized"
  }
  s.Set("userId", id)
  return 200, "user id is " + id
}

func Home(r render.Render) {
 
    news, err := news.GetLastNews()
    if err != nil {
      fmt.Println(err.Error())
    }
    r.HTML(200, "home", news)
}

func Kotoba(r render.Render, u *User) {
 
    k, err := kotoba.GetAllKotoba(u.Id)
    if err != nil {
      fmt.Println(err.Error())
    }
    r.HTML(200, "kotobalist", k)
	}

func KotobaId(params martini.Params, r render.Render) {
    id, _ := strconv.Atoi(params["id"])
    k, err := kotoba.GetKotoba(id)
    if err != nil {
      fmt.Println(err.Error())
    }
		r.HTML(200, "kotoba", k)
}
