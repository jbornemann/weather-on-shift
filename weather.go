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
    "strconv"
)

var (
    //This environment variable is fed from the downward API. See the deployment configuration for details
    ENVIRONMENT string = os.Getenv("NAMESPACE")
)


type WeatherSession struct {
    WeatherAPIKey string
    CurrentTemp   float64
    FetchIntervalMinutes uint32
}

func NewWeatherSession() (*WeatherSession) {
    log.Println("Creating new weather session")
    weather_session := &WeatherSession{}
    log.Println("Fetching key..")
    weather_session.fetch_key()
    log.Println("Updating weather store..")
    weather_session.update_weather()
    fetch_update_env := os.Getenv("UPDATE_INTERVAL_MINUTES")
    if fetch_update, err := strconv.Atoi(fetch_update_env); err == nil && fetch_update > 0 {
        weather_session.FetchIntervalMinutes = uint32(fetch_update)
        log.Println("Setting update interval to ", fetch_update_env, " minutes")
    } else {
        weather_session.FetchIntervalMinutes = 60
        log.Println("No update interval set. Defaulting to an hour")
    }
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
    //Note we use the external service endpoint here in the form open-weathermap-service.<namespace>.svc
    resp, err := http.Get(fmt.Sprintf("http://open-weathermap-service.%s.svc/data/2.5/weather?q=Raleigh&units=imperial&appid=%s", ENVIRONMENT, session.WeatherAPIKey))
    if err != nil || resp.StatusCode != 200 {
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
        sleep_time := time.Duration(session.FetchIntervalMinutes) * time.Minute
        time.Sleep(sleep_time)
        session.update_weather()
    }
}

func main() {
    weather_session := NewWeatherSession()
    go update_timer(weather_session)
    serve_weather(weather_session)
}
