package util

import (
  "testing"
)


func Test_GetKotoba(t *testing.T) {
  user := "serash"
  t.Log("get kotoba for " + user)
  getAllKotoba(user)  
}