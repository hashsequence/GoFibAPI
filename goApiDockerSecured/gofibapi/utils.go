package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"os"
	"time"
  "runtime"
)

//configuration variable handles conifgration for server
var config Configuration
//logger
var logger *log.Logger



func NewServer() *FibServer {
	//using a multiplexer to handle http requests
  mux := http.NewServeMux()
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	return &FibServer{server, mux}
}

//Endpoints are located here
func (s *FibServer) Routes() {
	 c := NewCurrentHandler()
	 n := NewNextHandler()
	 p := NewPreviousHandler()
	 //handles default endpoint /
	s.mux.HandleFunc("/", index)
	//handles current endpoint and wrapping it in a recovery function to handle panics(errors)
	s.mux.Handle("/current", RecoverWrap(c))
	//handles next endpoint
	s.mux.Handle("/next", RecoverWrap(n))
	//handles preious endpoints
	s.mux.Handle("/previous", RecoverWrap(p))
}

//wrapper function for http.server.ListenAndServe
func (s *FibServer) ListenAndServe() {
		s.Server.ListenAndServe()
}


func init() {
	///instantiates global variable to hold fibonaccie sequence, it is concurrent and thread safe
	FibArr = NewConcurrArr()
	FibArr.Set(0,NewLargeInt("nil"))
	FibArr.Set(1,NewLargeInt("nil"))
	FibArr.Set(2,NewLargeInt("0"))
	//loads config file from config.json with address
	loadConfig()
	//create log file
	file, err := os.OpenFile("fibApi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
}

//loads config.json file in root
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}

// for logging
func info(args ...interface{}) {
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

func danger(args ...interface{}) {
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}


func (s *FibServer	) initMsg() {
	fmt.Println("fibApi started at", ": " + config.Address)
}

//prints memory usage for testing purposes
func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
