package main

func main() {
  //inits the server and configuration
  fibServer := NewServer()
  //prints to console a message of server starting
  fibServer.initMsg()
  //all handlers for endpoints are run in Routes
  fibServer.Routes()
  //listens for http requests and responds 
  fibServer.ListenAndServe()
}
