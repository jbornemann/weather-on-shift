package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
    "io/ioutil"
    "strings"
    "github.com/tidwall/gjson"
    "time"
)

type WeatherSession struct {
    WeatherAPIKey string
    CurrentTemp   float64
}

func NewWeatherSession() (*WeatherSession) {
    weather_session := &WeatherSession{}
    weather_session.fetch_key()
    weather_session.update_weather()
    return weather_session
}

func (session *WeatherSession) fetch_key() {
    const KEY_FILE string = "/data/weather/key"
    log.Println("Reading weather key file..")
    key, err := ioutil.ReadFile(KEY_FILE)
    if err != nil {
        log.Fatal(fmt.Sprint("Could not read the weather api key from the secret", KEY_FILE))
        os.Exit(1)
    }
    session.WeatherAPIKey = strings.TrimSpace(string(key))
}

func (session *WeatherSession) update_weather() {
    log.Println("Requesting current weather..")
    resp, err := http.Get(fmt.Sprint("http://open-weathermap-service.myproject.svc/data/2.5/weather?q=Raleigh&units=imperial&appid=", session.WeatherAPIKey))
    if err != nil {
        log.Fatal("Something went wrong when attempting to get the weather")
        os.Exit(1)
    }
    data, _ := ioutil.ReadAll(resp.Body)
    session.CurrentTemp = gjson.Get(string(data), "main.temp").Num
    log.Println("Weather updated to ", session.CurrentTemp)
}

func serve_weather(weather *WeatherSession) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println(fmt.Sprint("Received request from ", r.RemoteAddr, " at ", time.Now().String()))
        fmt.Fprintf(w, "Current Raleigh temperature: %f \n", weather.CurrentTemp)
    })
    log.Print("Ready for weather requests..")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func update_timer(session *WeatherSession) {
    for {
        time.Sleep(60 * time.Minute)
        session.update_weather()
    }
}

func main() {
    weather_session := NewWeatherSession()
    go update_timer(weather_session)
    serve_weather(weather_session)
}
