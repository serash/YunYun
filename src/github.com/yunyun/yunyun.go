// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "fmt"
  "github.com/yunyun/kotoba"
  "github.com/yunyun/news"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "strconv"
  //"time"
)


/*
 * Main function
 */
func main() {
  m := martini.Classic()
  //render html templates from directory
  m.Use(render.Renderer(render.Options{
    Layout: "layout",
  }))
  
  m.Get("/", func(r render.Render) {
 
    news, err := news.GetLastNews()
    if err != nil {
      fmt.Println(err.Error())
    }
    r.HTML(200, "home", news)
	})
  
  m.Get("/kotoba",func(r render.Render) {
 
    k, err := kotoba.GetAllKotoba(2)
    if err != nil {
      fmt.Println(err.Error())
    }
    r.HTML(200, "kotobalist", k)
	})
  
  
	m.Get("/kotoba/:id", func(params martini.Params, r render.Render) {
    id, _ := strconv.Atoi(params["id"])
    k, err := kotoba.GetKotoba(id)
    if err != nil {
      fmt.Println(err.Error())
    }
		r.HTML(200, "kotoba", k)
	})
 
	m.Run()
}
