package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
    "io/ioutil"
    "strings"
    "github.com/tidwall/gjson"
)

func main() {
    const KEY_FILE string = "/data/weather/key"
    log.Println("Reading weather key file..")
    key, err := ioutil.ReadFile(KEY_FILE)
    if err != nil {
        log.Fatal(fmt.Sprint("Could not read the weather api key from the secret", KEY_FILE))
        os.Exit(1)
    }
    san_key := strings.TrimSpace(string(key))
    log.Println("Requesting current weather..")
    resp, err := http.Get(fmt.Sprint("http://api.openweathermap.org/data/2.5/weather?q=Raleigh&units=imperial&appid=", san_key))
    if err != nil {
        log.Fatal("Something went wrong when attempting to get the weather")
        os.Exit(1)
    }
    data, _ := ioutil.ReadAll(resp.Body)
    temp := gjson.Get(string(data), "main.temp").Num
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Current Raleigh temperature: %f \n", temp)
    })
    log.Print("Ready for weather requests..")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
