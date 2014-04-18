package yunyun

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
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