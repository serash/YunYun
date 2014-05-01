package kotoba

import (
  "time"
  "strconv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

/*
 * Constants/Variables definitions
 */
// todo make a var table with level time
var time_per_level []string = []string{
  "1h", 
  "2h", 
  "4h", 
  "12h", 
  "24h", 
  "48h", 
  "192h", 
  "384h",
  "768h",  
  "1536h",
}
var insert_query string = "INSERT INTO " + 
                   "kotoba(user_id, kotoba, hatsuon, hatsuon_mnemonic, imi, "+
                   "imi_mnemonic, level, next_review, unlocked) " + 
                   "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
var update_review_query string = "UPDATE kotoba " + 
                   "SET level=?, " +
                   "next_review=? " +
                   "WHERE id = ?"
var select_all_query string = "SELECT * FROM kotoba WHERE (user_id = ?)"
var select_query string = "SELECT * FROM kotoba WHERE (id = ?)"
var select_review_query string = "SELECT * FROM kotoba WHERE "+
                        "(user_id = ? AND next_review < NOW())"
var select_first_review_query string = "SELECT * FROM kotoba WHERE "+
                        "(user_id = ? AND next_review < NOW()) " + 
                        "ORDER BY unlocked DESC LIMIT 1"
var select_next_review_query string = "SELECT * FROM kotoba WHERE "+
                        "user_id = ? ORDER BY next_review ASC LIMIT 1"

type Kotoba struct {
  Id        int
  User_id   int
  Goi    string
  Hatsuon   string
  Imi       string
  Hatsuon_  string // mnemonic for hatsuon
  Imi_      string // mnemonic for hatsuon
  Level     int
  Review    string
  Unlocked  string
  Valid     int
}
func FormatTime(t time.Time) string {
  return t.Format("2006-01-02 15:04:05")
}
func ExistingKotoba(id_ int, uid int, k string, h string, i string, h_ string, i_ string, 
                    l int, r string, u string) *Kotoba {
  return &Kotoba{Id: id_, User_id:uid, Goi: k, Hatsuon: h, Imi: i, Hatsuon_ : h_, 
                 Imi_: i_, Level: l, Review: r, Unlocked: u}
}
func NewKotoba(uid int, k string, h string, i string, l int, r string) *Kotoba {
  return &Kotoba{Id: -1, User_id:uid, Goi: k, Hatsuon: h, Imi: i, Level: l, 
                 Review: r, Unlocked: FormatTime(time.Now().Local())}
}
func NewDefaultLevelKotoba(uid int, k string, h string, h_ string, i string, i_ string) *Kotoba {
  timenow := time.Now().Local()
  hours, _ := time.ParseDuration(time_per_level[2])
  r := timenow.Add(hours)
  return &Kotoba{Id: -1, User_id:uid, Goi: k, Hatsuon: h, Hatsuon_: h_, Imi: i, Imi_: i_, Level: 2, 
                 Review: FormatTime(r), Unlocked: FormatTime(time.Now().Local())}
}

/* 
 * kotoba functions
 */
func (k *Kotoba) IncLevel() {
  k.Level++
  if (k.Level > 9) {
    k.Level = 9
  }
  hours, _ := time.ParseDuration(time_per_level[k.Level])    
  r, _ := time.Parse("2006-01-02 15:04:05", k.Review)
  k.Review = FormatTime(r.Add(hours))
}
func (k *Kotoba) DecLevel() {
  k.Level--
  if (k.Level < 0) {
    k.Level = 0
  }
  hours, _ := time.ParseDuration(time_per_level[k.Level]) 
  r, _ := time.Parse("2006-01-02 15:04:05", k.Review)
  k.Review = FormatTime(r.Add(hours))
}
func (k *Kotoba)Save(db *sql.DB) (error) {
  // insert kotoba
  _, err := db.Exec(insert_query,
                   k.User_id, k.Goi, k.Hatsuon, k.Hatsuon_, k.Imi, 
                   k.Imi_, k.Level, k.Review, k.Unlocked)
  return err
}
func (k *Kotoba)Update(db *sql.DB) (error) {
  // insert kotoba
  var err error
  if(k.Id == -1) {
    _, err = db.Exec(insert_query,
                   k.User_id, k.Goi, k.Hatsuon, k.Hatsuon_, k.Imi, 
                   k.Imi_, k.Level, k.Review, k.Unlocked)
  } else {
    _, err = db.Exec(update_review_query, k.Level, k.Review, k.Id)
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
    if err != nil {
      return -1
    }
    return v
  }
}
func GetKotobaFromRow(row *sql.Row) *Kotoba {
  // get data from row
  var k Kotoba
  var h sql.NullString
  var h_ sql.NullString
  var i sql.NullString
  var i_ sql.NullString
  err := row.Scan(&k.Id, &k.User_id, &k.Goi, &h, &h_, &i, &i_, &k.Level, &k.Review, &k.Unlocked)
  if err != nil {
    k.Id = -1
    k.Valid = -1
    return &k
  }
  k.Hatsuon = ""
  k.Hatsuon_ = ""
  k.Imi = ""
  k.Imi_ = ""
  if h.Valid { 
    k.Hatsuon = h.String
  }
  if h_.Valid {
    k.Hatsuon_ = h_.String
  }
  if i.Valid {
    k.Imi = i.String
  }
  if i_.Valid {
    k.Imi_ = i_.String
  }
  k.Valid = 0
  return &k
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
func compare(answer string, ref string) (bool) {
  return answer == ref
}

func GetKotoba(id int, db *sql.DB) (*Kotoba) {
  // query kotoba
  row := db.QueryRow(select_query, id)
  // get data from rows
  kotoba := GetKotobaFromRow(row)
  return kotoba
}
func GetAllKotoba(user_id int, db *sql.DB) (*[]Kotoba, error) {
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
func GetReviewKotoba(user_id int, db *sql.DB) (*[]Kotoba, error) {
  // query kotoba
  rows, err := db.Query(select_review_query, user_id)
  if err != nil {
    panic(err.Error())
  }
  // get data from rows
  defer rows.Close()
  kotobaArray := GetKotobaFromRows(rows)
  return kotobaArray, err
}
func GetFirstReviewKotoba(user_id int, db *sql.DB) (*Kotoba) {
  // query kotoba
  row := db.QueryRow(select_first_review_query, user_id)
  kotoba := GetKotobaFromRow(row)
  if kotoba.Valid == -1 {
    row := db.QueryRow(select_next_review_query, user_id)
    kotoba = GetKotobaFromRow(row)
    kotoba.Valid = -1
  }
  return kotoba
}