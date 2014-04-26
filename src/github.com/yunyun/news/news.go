package news

import (
  "fmt"
  "time"
  "strconv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/yunyun/db"
)

var insert_query string = "INSERT INTO " + 
                   "news(title, description)"+
                   "VALUES(?, ?)"
var update_query string = "UPDATE news " + 
                   "SET title=?, " +
                   "SET description=?, " +
                   "WHERE id = ?"
var select_all_query string = "SELECT * FROM news"
var select_5_last_query string = "SELECT * FROM news ORDER BY date DESC LIMIT 5"

type News struct {
  Id            int
  Title         string
  Date          string
  Description   string
}
func FormatTime(t time.Time) string {
  //return t.Format("2006-01-02")
  return t.Format("02 January '06")
}

/* 
 * news functions
 */
func (n *News)SaveNews() (error) {
  // get database
  db, err := db.GetDB()
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()
  // insert kotoba
  _, err = db.Exec(insert_query, n.Title, n.Description)
  return err
}
func (n *News)UpdateNews() (error) {
  // get database
  db, err := db.GetDB()
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()
  // insert kotoba
  if(n.Id == -1) {
    _, err = db.Exec(insert_query, n.Title, n.Description)
  } else {
    _, err = db.Exec(update_query, n.Title, n.Description, n.Id)
  }
  return err
}
func GetStringValue(val []byte) string {
  if val == nil {
    return "NULL"
  } else {
    return string(val)
  }
}
func GetIntValue(val []byte) int {
  if val == nil {
    return -1
  } else {
    v, err := strconv.Atoi(string(val))
    if err == nil {
      return -1
    }
    return v
  }
}
func GetNewsFromRow(rows *sql.Rows) *[]News {
  // get data from rows
  var newsArray []News
  columns, err := rows.Columns()
  values := make([]sql.RawBytes, len(columns))
  scanArgs := make([]interface{}, len(values))
  for i := range values {
    scanArgs[i] = &values[i]
  }
  //var news News
  for rows.Next() {
    err = rows.Scan(scanArgs...)
    if err != nil {
      panic(err.Error())
    }
    var news News
    news.Id = GetIntValue(values[0])
    news.Title = GetStringValue(values[1])
    news.Date = GetStringValue(values[2])
    date, _ := time.Parse("2006-01-02 15:04:05", news.Date)
    news.Date = FormatTime(date)
    news.Description = GetStringValue(values[3])
    newsArray = append(newsArray, news)
  }
  fmt.Println()
  return &newsArray
}
func GetAllNews() (*[]News, error) {
  // get database
  db, err := db.GetDB()
  if err != nil {
     panic(err.Error())
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
  }
  // query kotoba
  rows, err := db.Query(select_all_query)
  if err != nil {
    panic(err.Error())
  }  
  // get data from rows
  defer rows.Close()
  newsArray := GetNewsFromRow(rows)
  return newsArray, err
}
func GetLastNews() (*[]News, error) {
  // get database
  db, err := db.GetDB()
  if err != nil {
     panic(err.Error())
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
      panic(err.Error()) // proper error handling instead of panic in your app
  }
  // query kotoba
  rows, err := db.Query(select_5_last_query)
  if err != nil {
    panic(err.Error())
  }
  // get data from rows
  defer rows.Close()
  newsArray := GetNewsFromRow(rows)
  return newsArray, err
}