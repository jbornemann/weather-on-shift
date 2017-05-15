package main

import (
    "fmt"
    "net/http"
    "log"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello world!\n")
    })
    log.Print("Ready for weather requests..")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
