package main

import (
  "fmt"
  "encoding/json"
//  "log"
  "net/http"
  "errors"
)

//I created my own data structure to hold my integers, which consists of a string, since a string can hold more digits than int64
//max length of value that fibonacci number will go to
const MAX_LEN = 100
//concurrent and thread safe data store to store current, previous, and previous's previous fibonacci sequence
var  FibArr *ConcurrArrOfLargeInt

//handles / endpoint, and prints memory usage
func index(w http.ResponseWriter, request *http.Request) {
  fmt.Println("this is the index")
   PrintMemUsage()
}

//handler function for current, writes the current fibonacci sequence number into a json as a string as a response
func (f *CurrentHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {

  res := FibNum{FibArr.Get(2).GetVal()}
  currJson, err := json.MarshalIndent(res, "", "  ")
  if err != nil {
    danger("route_main.go : error : func : Current: ", err)
    //log.Fatal("route_main.go : error : func : Current: ", err)
  }
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  w.Write(currJson)

}

//handler function for next, shifts the sequence forward and responds with the current fib sequence in a json holding the value in a string
func (f *NextHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {

  num := FibArr.ShiftForward()
  res := FibNum{num}
  nextJson, err := json.MarshalIndent(res, "", "  ")
  if err != nil {
    danger("route_main.go : error : func : Next: ", err)
  //  log.Fatal("route_main.go : error : func : Next: ", err)
  }
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  w.Write(nextJson)
}

//handler function for previous, shifts the sequence backward and responds with the current fib sequence in a json holding the value in a string
func (f *PreviousHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {

  num := FibArr.ShiftBackward()
  res := FibNum{num}
  previousJson, err := json.MarshalIndent(res, "", "  ")
  if err != nil {
    danger("route_main.go : error : func : Previous: ", err)
    //log.Fatal("route_main.go : error : func : Previous: ", err)
  }
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  w.Write(previousJson)
}


//recovery wrapper for handlers and will handle crashes from handlers
func RecoverWrap(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var err error
        defer func() {
            r := recover()
            if r != nil {
                switch t := r.(type) {
                case string:
                    err = errors.New(t)
                case error:
                    err = t
                default:
                    err = errors.New("Unknown error")
                }
                danger("App Crashed ... Recovering from ",err.Error())
                FibArr.Reset()
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        }()
        h.ServeHTTP(w, r)
    })
}
