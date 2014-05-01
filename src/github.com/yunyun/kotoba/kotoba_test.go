package kotoba

import (
  "fmt"
  "testing"
  "github.com/yunyun/db"
)


func Test_GetKotoba(t *testing.T) {  
  db, err := db.GetDB()
  if err != nil {
    panic(err)
  }
  k := GetKotoba(6, db)
  fmt.Println(k)
  k.IncLevel()
  fmt.Println(k)
}
