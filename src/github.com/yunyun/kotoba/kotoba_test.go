package kotoba

import (
  "fmt"
  "testing"
  "github.com/yunyun/db"
)


func Test_MongoDB(t *testing.T) {  
  db, err := db.GetDB()
  if err != nil {
    panic(err)
  }
  k := GetKotoba(66, db)
  fmt.Println(k)
  k.IncLevel()
  fmt.Println(k)
  fmt.Println(TimeUntil(k.Review))
  ks, _ := FindKotoba(2, "america", db)
  fmt.Println(ks)
}
