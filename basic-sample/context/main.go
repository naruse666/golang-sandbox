package main

import (
  "fmt"
  "net/http"
  "time"
)

func hello(w http.ResponseWriter, req *http.Request) {
  ctx := req.Context()
  fmt.Println("started server")
  defer fmt.Println("end server")

  select {
  case <-time.After(10 * time.Second) :
    fmt.Fprintf(w, "hello \n")
  case <-ctx.Done():
    err := ctx.Err()
    fmt.Println("Server:", err)
    internalError := http.StatusInternalServerError
    http.Error(w, err.Error(), internalError)
  }
}

func main() {
  http.HandleFunc("/hello", hello)
  http.ListenAndServe(":8090", nil)
}
