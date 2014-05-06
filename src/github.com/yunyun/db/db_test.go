package db

import (
  "fmt"
  "testing"
  "strconv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "labix.org/v2/mgo/bson"
)

func Test_GetReviews(t *testing.T) { 
  //c, err := GetKotobaCollection()
  //IfPanic(err)
  //db_, err := GetMongoDB()
  //IfPanic(err)
  // find user
  //user, err := GetUser("serash", db_)
  //fmt.Println(user)
  // find all words of user
  //kotoba := GetAllKotoba(user.Id, c)
  //fmt.Println(kotoba)
  // find new reviews
  //c, _ = cKotoba.Find(bson.M{"user_id": user.Id, "review": bson.M{"$lte": time.Now()}}).Count()
  //fmt.Println(c)
}
func GetSQLUser(user int, db *sql.DB) (*MongoUser, error) {
  u := MongoUser{}
  var t string
  err := db.QueryRow("SELECT user, email, pw_hash, unlocked from users where id=?",
                     user).Scan(&u.Name, &u.Mail, &u.Pw_hash, &t)
  u.Unlocked = StringToTime(t)
  u.Id = bson.NewObjectId()
  fmt.Println(u)
  return &u, err
}
var select_all_query string = "SELECT * FROM kotoba WHERE (user_id = ?) ORDER BY next_review ASC"
type Kotoba struct {
  Id        int
  User_id   int
  Goi       string
  Hatsuon   string
  Imi       string
  Hatsuon_  string // mnemonic for hatsuon
  Imi_      string // mnemonic for hatsuon
  Level     int
  Review    string
  ReviewUntil string
  Unlocked  string
  Valid     int
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
    if err != nil {
      return -1
    }
    return v
  }
}
func GetKotobaFromRows(rows *sql.Rows) *[]Kotoba {
  // get data from rows
  var kotobaArray []Kotoba
  columns, err := rows.Columns()
  values := make([]sql.RawBytes, len(columns))
  scanArgs := make([]interface{}, len(values))
  for i := range values {
    scanArgs[i] = &values[i]
  }
  //var kotoba Kotoba
  for rows.Next() {
    err = rows.Scan(scanArgs...)
    if err != nil {
      panic(err.Error())
    }
    var kotoba Kotoba
    kotoba.Id = GetIntValue(values[0])
    kotoba.User_id = GetIntValue(values[1])
    kotoba.Goi = GetStringValue(values[2])
    kotoba.Hatsuon = GetStringValue(values[3])
    kotoba.Hatsuon_ = GetStringValue(values[4])
    kotoba.Imi = GetStringValue(values[5])
    kotoba.Imi_ = GetStringValue(values[6])
    kotoba.Level = GetIntValue(values[7])
    kotoba.Review = GetStringValue(values[8])
    kotoba.Unlocked = GetStringValue(values[9])
    //fmt.Println(kotoba)
    kotobaArray = append(kotobaArray, kotoba)
  }
  return &kotobaArray
}
func GetAllSQLKotoba(user_id int, db *sql.DB) (*[]Kotoba, error) {
  // query kotoba
  rows, err := db.Query(select_all_query, user_id)
  if err != nil {
    panic(err.Error())
  }  
  // get data from rows
  defer rows.Close()
  kotobaArray := GetKotobaFromRows(rows)
  return kotobaArray, err
}
func Test_MysqlToMongo(t *testing.T) {
  db_, err := GetMongoDB()
  cUser := db_.C(MDB_USERS)
  IfPanic(err)
  cKotoba, err := GetKotobaCollection()
  IfPanic(err)
  db, err := GetDB()
  IfPanic(err)
  u1, err := GetSQLUser(2, db)
  IfPanic(err)
  u2, err := GetSQLUser(3, db)
  IfPanic(err)
  err = cUser.Insert(u1)
  IfPanic(err)
  err = cUser.Insert(u2)
  IfPanic(err)
  serash := MongoUser{}
  cUser.Find(bson.M{"name" : "serash"}).One(&serash)
  fmt.Println(serash)
  
  k, err := GetAllSQLKotoba(2, db)
  for _, k_ := range *k {
    newK := MongoKotoba{Id: bson.NewObjectId(), 
                      User_id: serash.Id, 
                      Goi: k_.Goi, 
                      Hatsuons: []HatsuonKotoba{HatsuonKotoba{Hatsuon: k_.Hatsuon} }, 
                      Hatsuon_: k_.Hatsuon_, 
                      Imis: []ImiKotoba{ImiKotoba{Imi: k_.Imi} }, 
                      Imi_: k_.Imi_, 
                      Level: k_.Level, 
                      Review: StringToTime(k_.Review),
                      Unlocked: StringToTime(k_.Unlocked)}
    fmt.Println(newK)
    err = cKotoba.Insert(&newK)
    IfPanic(err)
  }
}
