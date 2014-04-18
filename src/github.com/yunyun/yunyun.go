// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "github.com/yunyun/utils"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)
var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
)

/* 
 * test functions
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
	http.HandleFunc("", RedirectHandler)
	http.HandleFunc("/account", AccountHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)

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
