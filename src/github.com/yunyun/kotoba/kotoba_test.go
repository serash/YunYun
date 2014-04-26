package kotoba

import (
  "time"
  "fmt"
  "testing"
)


func Test_GetKotoba(t *testing.T) {
  user_id := 2 
  
  k, err := GetAllKotoba(user_id)
  if err != nil {
    fmt.Println(err.Error())
  }
  for _, s := range *k {
    fmt.Println(s)
  }
  time := time.Now()
  fmt.Println(time.Format("2006-01-02 15:04:05"))
  ko, err:= GetKotoba(2)
  fmt.Println(ko)
  
}
