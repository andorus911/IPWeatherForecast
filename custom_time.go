package main

import (
	"time"
  "strings"
  "fmt"
)

type cTime struct {
	time.Time
}

var layout string = "2006-01-02T15:04"

func (ct *cTime) UnmarshalJSON(data []byte) (err error) {
  s := strings.Trim(string(data), "\"")
  if s == "null" {
     ct.Time = time.Time{}
     return
  }
  
	parsedTime, err := time.Parse(layout, s)
	if err != nil {
		return
	}
  
	ct.Time = parsedTime
	return
}

func (ct *cTime) MarshalJSON() ([]byte, error) {
  if ct.Time.IsZero() {
    return []byte("null"), nil
  }
  return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(layout))), nil
}