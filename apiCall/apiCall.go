package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
  "time"
)

type Num struct {
	CurrNum  string
}

func main() {

  start := time.Now()
	url := "http://0.0.0.0:8080/next"
  for i := 0; i <  1000; i++ {
      apiCall(url)
  }


    elapsed := time.Since(start)
  fmt.Println(elapsed)



}

func apiCall(url string) {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal("NewRequest: ", err)
    return
  }

  // For control over HTTP client headers,
  // redirect policy, and other settings,
  // create a Client
  // A Client is an HTTP client
  client := &http.Client{}

  // Send the request via a client
  // Do sends an HTTP request and
  // returns an HTTP response
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Do: ", err)
    return
  }

  // Callers should close resp.Body
  // when done reading from it
  // Defer the closing of the body
  defer resp.Body.Close()

  // Fill the record with the data from the JSON
  var record Num

  // Use json.Decode for reading streams of JSON data
//  fmt.Println(resp.Body)
  if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
    log.Println(err)
  }
  
  fmt.Println(record)
}
