package news

import (
  "fmt"
  "testing"
)


func Test_GetNews(t *testing.T) {
  
  n, err := GetLastNews()
  if err != nil {
    fmt.Println(err.Error())
  }
  for _, s := range *n {
    fmt.Println(s)
  }

}
