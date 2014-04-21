package util

import (
  "fmt"
  "log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

/*
 * Type definitions
 */
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
 * kotoba functions
 */
func getAllKotoba(user string) ([]Kotoba, error) {
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
  // query kotoba
  rows, err := db.Query("SELECT pw_hash FROM users WHERE (user = ? OR email = ?)", user)
  switch {
    case err == sql.ErrNoRows:
      log.Printf("No user found.")
    case err != nil:
            log.Fatal(err)
  }
  // get data from rows
  cols, _ := rows.Columns()
  n := len(cols)
  var kotobas []Kotoba
  // loop over
  var fields []interface{}
  for i := 0; i < n; i++ {
    fields = append(fields, new(string))
  }
  for rows.Next() {
    rows.Scan(fields...)
    //var ktb = Kotoba{kotoba: kotoba, imi: imi, ...}
    for i := 0; i < n; i++ {
      fmt.Print(*(fields[i].(*string)), "\t")
    }
    fmt.Println()
  }
  fmt.Println()
  return kotobas, err
}