// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "labix.org/v2/mgo"
  "github.com/yunyun/auth"
  "github.com/yunyun/kotoba"
  "github.com/yunyun/db"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"  
  "github.com/martini-contrib/sessions"
)
const (
  PAGE_ERROR string = "/error/"
  FUNC_SEARCH string = "/search"
  FUNC_HOME string = "/"
  PAGE_ADD string = "/addWord"
  PAGE_STATS string = "/stats"
  PAGE_EDIT string = "/edit/:id"
  ADD_TEMPLATE string = "add"
  FUNC_KOTOBA string = "/kotoba/:id"
  FUNC_ADD string = "/add"
  FUNC_EDIT string = "/edit"
  PAGE_REVIEW string = "/doReviews"
  FUNC_CHECK string = "/reviewed/:id"
  // todo later
  apiReviews string = "/apiReviews/:apikey"
)

func GetMongoDB() *mgo.Database {
  datab, err := db.GetMongoDB()
  db.IfPanic(err)
  return datab
}
func GetKotobaCollection() *mgo.Collection {
  c, err := db.GetKotobaCollection()
  db.IfPanic(err)
  return c
}

/*
 * Main function
 */
func main() {
  m := martini.Classic()
  
  //sessions
  store := sessions.NewCookieStore([]byte("yyt726ddd318"))
  m.Use(sessions.Sessions("yunyun", store))
  
  m.Map(GetMongoDB())
  m.Map(GetKotobaCollection())
  //render html templates from directory
  m.Use(render.Renderer(render.Options{
    Layout: "layout",
  }))
  
  m.Post(auth.FUNC_LOGIN, auth.PostLogin)  
  m.Post(auth.FUNC_REGISTER, auth.Register)
  m.Get(auth.FUNC_LOGOUT, auth.Logout)
  m.Post(FUNC_SEARCH, auth.RequireLogin, kotoba.Search)
  m.Get(FUNC_HOME, auth.RequireLogin, kotoba.Home)   
  m.Get(PAGE_ADD, auth.RequireLogin, AddWord) 
  m.Get(FUNC_KOTOBA, auth.RequireLogin, kotoba.ShowKotoba)  
  m.Post(FUNC_ADD, auth.RequireLogin, kotoba.AddKotoba)
  m.Get(PAGE_REVIEW, auth.RequireLogin, kotoba.DoReviews)
  m.Post(FUNC_CHECK, auth.RequireLogin, kotoba.CheckReview)
  m.Get(PAGE_EDIT, auth.RequireLogin, kotoba.EditKotoba) 
  m.Post(FUNC_EDIT, auth.RequireLogin, kotoba.SaveEditKotoba)
  m.Get(PAGE_STATS, auth.RequireLogin, kotoba.ShowStats)
  // todo: not yet implemented
  //m.Get(apiReviews, auth.RequireLogin, GetReviews)
 
	m.Run()
}

func AddWord(r render.Render) {
  r.HTML(200, ADD_TEMPLATE, nil)
}

