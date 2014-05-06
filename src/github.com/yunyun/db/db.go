package db

import (
  "fmt"
  "math/rand"
  "time"
  "strings"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)

/*
 * Constants definitions
 */
const (
  DB_HOST = "tcp(127.0.0.1:3306)"
  DB_NAME = "yunyun"
  DB_USER = "yunyun"
  DB_PASS = "nuynuy"
  MDB_HOST = "localhost"
  MDB_DB = "yunyun"
  MDB_USERS = "users"
  MDB_KOTOBA = "kotoba"
  DEFAULT_REVIEW_TIM = "4h"
)
var time_per_level []string = []string{
  "1h", // beginner
  "2h", // beginner
  "4h", // beginner
  "12h", // elementary
  "24h", //elementary
  "48h", // intermediate
  "192h", // intermediate
  "384h", // master
  "768h",  // master
  "1536h", // known
}
/*
 * structs
 */

type MongoUser struct {
  Id bson.ObjectId `bson:"_id"`
  Name string      `bson:"name"`
  Mail string      `bson:"mail"`
  Pw_hash string    `bson:"pw_hash"`
  Unlocked time.Time  `bson:"unlocked"`
}
type LabelKotoba struct {
  Label string `bson:"label"`
}
type ImiKotoba struct {
  Imi string `bson:"imi"`
}
type HatsuonKotoba struct {
  Hatsuon string `bson:"hatsuon"`
}
type MongoKotoba struct {
  Id bson.ObjectId         `bson:"_id"`
  User_id bson.ObjectId    `bson:"user_id"`
  Goi string               `bson:"goi"`
  Hatsuons []HatsuonKotoba `bson:"hatsuons"`
  Hatsuon_ string          `bson:"hatsuon_"`
  Imis []ImiKotoba         `bson:"imis"`
  Imi_ string              `bson:"imi_"`
  Labels []LabelKotoba     `bson:"labels"`
  Level int                `bson:"level"`
  Difficulty int           `bson:"difficulty"`
  Review time.Time         `bson:"review"`
  Unlocked time.Time       `bson:"unlocked"`
}
func IfPanic(err error) {
  if err != nil {
    panic(err)
  }
}
func TimeToString(t time.Time) string {
  return t.Format("2006-01-02 15:04:05")
}
func StringToTime(t string) time.Time {
  tt, _ := time.Parse("2006-01-02 15:04:05", t)  
  return tt
}
/* 
 * database functions
 */
func GetMongoDB() (*mgo.Database, error) {
  sess, err := mgo.Dial(MDB_HOST)
  return sess.DB(MDB_DB), err
}
func GetKotobaCollection() (*mgo.Collection, error) {
  sess, err := mgo.Dial(MDB_HOST)
  return sess.DB(MDB_DB).C(MDB_KOTOBA), err
}
// USERS
func GetUserById(id bson.ObjectId, db *mgo.Database) (*MongoUser, error) {
  user := MongoUser{}
  err := db.C(MDB_USERS).Find(bson.M{"_id": id }).One(&user)
  return &user, err
}
func GetUser(name string, db *mgo.Database) (*MongoUser, error) {
  user := MongoUser{}
  err := db.C(MDB_USERS).Find(bson.M{"$or": []interface{}{bson.M{"name" : name}, bson.M{"mail" : name} } } ).One(&user)
  return &user, err
}
func InsertUser(name string, mail string, pw_hash string, db *mgo.Database) error {
  user := MongoUser{Id : bson.NewObjectId(), Name : name, Mail : mail, Pw_hash : pw_hash, Unlocked: time.Now()}
  return db.C(MDB_USERS).Insert(&user)
}
// KOTOBA
func GetKotoba(id bson.ObjectId, c *mgo.Collection) (*MongoKotoba) {
  kotoba := MongoKotoba{}
  c.Find(bson.M{"_id": id}).One(&kotoba)
  return &kotoba
}
func GetAllKotoba(user_id bson.ObjectId, c *mgo.Collection) (*[]MongoKotoba) {
  kotoba := []MongoKotoba{}
  c.Find(bson.M{"user_id": user_id}).All(&kotoba)
  return &kotoba
}
func FindKotoba(user_id bson.ObjectId, search string, c *mgo.Collection) (*[]MongoKotoba) {
  kotoba := []MongoKotoba{}
  c.Find(bson.M{"user_id": user_id, "$or": 
                []interface{}{bson.M{"goi" : bson.M{"$regex" : ".*" + search + ".*"}},
                              bson.M{"imis.imi" : bson.M{"$regex" : ".*" + search + ".*"}},
                              bson.M{"hatsuons.hatsuon" :  bson.M{"$regex" : ".*" + search + ".*"}}  }}).All(&kotoba)
  return &kotoba
}
func GetNumberBeginner(user_id bson.ObjectId, c *mgo.Collection) int {
  n, err := c.Find(bson.M{"user_id": user_id, "$or": 
                []interface{}{bson.M{"level" : 0},
                              bson.M{"level" : 1},
                              bson.M{"level" : 2} }}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberElementary(user_id bson.ObjectId, c *mgo.Collection) int {
  n, err := c.Find(bson.M{"user_id": user_id, "$or": 
                []interface{}{bson.M{"level" : 3},
                              bson.M{"level" : 4} }}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberIntermediate(user_id bson.ObjectId, c *mgo.Collection) int {
  n, err := c.Find(bson.M{"user_id": user_id, "$or": 
                []interface{}{bson.M{"level" : 5},
                              bson.M{"level" : 6} }}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberMaster(user_id bson.ObjectId, c *mgo.Collection) int {
  n, err := c.Find(bson.M{"user_id": user_id, "$or": 
                []interface{}{bson.M{"level" : 7},
                              bson.M{"level" : 8} }}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberKnown(user_id bson.ObjectId, c *mgo.Collection) int {
  n, err := c.Find(bson.M{"user_id": user_id, "level" : 9}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberReviewsNow(user_id bson.ObjectId, c *mgo.Collection) int {
  n, err := c.Find(bson.M{"user_id": user_id, 
                          "review": bson.M{"$lte": time.Now()}}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberReviewsHour(user_id bson.ObjectId, c *mgo.Collection) int {
  hours, _ := time.ParseDuration("1h")   
  t := time.Now().Local().Add(hours)
  n, err := c.Find(bson.M{"user_id": user_id, 
                          "review": bson.M{"$lte": t}}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetNumberReviewsDay(user_id bson.ObjectId, c *mgo.Collection) int {
  hours, _ := time.ParseDuration("24h")   
  t := time.Now().Local().Add(hours)
  n, err := c.Find(bson.M{"user_id": user_id, 
                          "review": bson.M{"$lte": t}}).Count()
  if err != nil {
    return 0
  }
  return n
}
func GetReviewKotoba(user_id bson.ObjectId, n int, c *mgo.Collection) (*[]MongoKotoba) {
  kotoba := []MongoKotoba{}
  c.Find(bson.M{"user_id": user_id, "review": bson.M{"$lte": time.Now()}}).Limit(1).All(&kotoba)
  return &kotoba
}
func GetNextKotoba(user_id bson.ObjectId, c *mgo.Collection) (*MongoKotoba) {
  kotoba := MongoKotoba{}  
  c.Find(bson.M{"user_id": user_id, "review": bson.M{"$gte": time.Now()}}).Sort("review").One(&kotoba)
  return &kotoba
}
func GetRandomKotoba(user_id bson.ObjectId, c *mgo.Collection) (*MongoKotoba) {
  kotoba := MongoKotoba{}  
  query := c.Find(bson.M{"user_id": user_id, "review": bson.M{"$lte": time.Now()}})
  n, _ := query.Count()
  if n > 0 {
    query.Skip(rand.Intn(n)).One(&kotoba)
    fmt.Println(kotoba)
    return &kotoba
  } else {
    return nil
  }
}
func (k *MongoKotoba) IncLevel() {
  k.Level++
  if (k.Level > 9) {
    k.Level = 9
  }
  hours, _ := time.ParseDuration(time_per_level[k.Level])   
  t := time.Now().Local()
  k.Review = t.Add(hours)
}
func (k *MongoKotoba) DecLevel() {
  k.Level--
  if (k.Level < 0) {
    k.Level = 0
  }
  hours, _ := time.ParseDuration(time_per_level[k.Level]) 
  t := time.Now().Local()
  k.Review = t.Add(hours)
}
func ReviewUpdateKotoba(id string, review bool, c *mgo.Collection) (error) {
  _id := bson.ObjectIdHex(id)
  kotoba := MongoKotoba{}
  c.Find(bson.M{"_id": _id}).One(&kotoba)
  if review {
    kotoba.IncLevel()
  } else {
    kotoba.DecLevel()
  }
  fmt.Println(kotoba)
  return c.UpdateId(_id, kotoba);
}
func SaveKotoba(mk *MongoKotoba, c *mgo.Collection) (error) {
  return c.UpdateId(mk.Id, mk);
}

func SplitLabels(val string) (*[]LabelKotoba) {
  vals := strings.Split(val, ",")
  arr := []LabelKotoba{}
  for _, v := range vals {
    newv := strings.TrimSpace(v)
    if len(newv) != 0 {    
      arr = append(arr, LabelKotoba{Label: strings.ToLower(newv)})
    }
  }  
  fmt.Println(arr)
  return &arr
}
func SplitImi(val string) (*[]ImiKotoba) {
  vals := strings.Split(val, ",")
  arr := []ImiKotoba{}
  for _, v := range vals {
    newv := strings.TrimSpace(v)
    if len(newv) != 0 {    
      arr = append(arr, ImiKotoba{Imi: strings.ToLower(newv)})
    }
  }  
  fmt.Println(arr)
  return &arr
}
func SplitHatsuon(val string) (*[]HatsuonKotoba) {
  vals := strings.Split(val, ",")
  arr := []HatsuonKotoba{}
  for _, v := range vals {
    newv := strings.TrimSpace(v)
    if len(newv) != 0 {
      arr = append(arr, HatsuonKotoba{Hatsuon: strings.ToLower(newv)})
    }
  }  
  fmt.Println(arr)
  return &arr
}
func SaveNewKotoba(user_id bson.ObjectId, word string, ll string, diff int, hatsuon string, h_ string,
                   imi string, i_ string, c *mgo.Collection) (*MongoKotoba, error) {
  timenow := time.Now().Local()
  hours, _ := time.ParseDuration(DEFAULT_REVIEW_TIM)
  r := timenow.Add(hours)
  imis := SplitImi(imi)
  hatsuons := SplitHatsuon(hatsuon)
  labels := SplitLabels(ll)
// add meanings
  newK := MongoKotoba{Id: bson.NewObjectId(), 
                      User_id: user_id, 
                      Goi: word, 
                      Hatsuons: *hatsuons, 
                      Hatsuon_: h_, 
                      Imis: *imis, 
                      Imi_: i_, 
                      Labels: *labels,
                      Level: 2, 
                      Review: r,
                      Difficulty: diff,
                      Unlocked: time.Now().Local()}
  return &newK, c.Insert(&newK)
}
func SaveEditKotoba(id string, word string, ll string, diff int, hatsuon string, h_ string,
                   imi string, i_ string, c *mgo.Collection) (*MongoKotoba, error) {
  imis := SplitImi(imi)
  hatsuons := SplitHatsuon(hatsuon)
  labels := SplitLabels(ll)
  _id := bson.ObjectIdHex(id)
  kotoba := MongoKotoba{}
  c.Find(bson.M{"_id": _id}).One(&kotoba)
  fmt.Println(kotoba)
  kotoba.Goi = word
  kotoba.Hatsuons = *hatsuons
  kotoba.Imis = *imis
  kotoba.Labels = *labels
  kotoba.Imi_ = i_
  kotoba.Hatsuon_ = h_
  kotoba.Difficulty = diff  
  fmt.Println(kotoba)
  return &kotoba, c.UpdateId(kotoba.Id, &kotoba)
}