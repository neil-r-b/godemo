package weatherapis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

var errorLogger = log.New(os.Stderr,
	"ERROR: ",
	log.Ldate|log.Ltime|log.Lshortfile)

// Interface for all APIs to get Weather
type WeatherAPI interface {
	GetWeather(zipcode string, ch chan WeatherResult)
}

// Contains the results from API calls
type WeatherResult struct {
	name string
	temp float64
	err  error
}

func (wr WeatherResult) GetTemp() float64 {
	return wr.temp
}

func (wr WeatherResult) GetError() error {
	return wr.err
}

func (wr WeatherResult) String() string {
	if wr.err == nil {
		return fmt.Sprintf("%s: %.2f℉", wr.name, wr.temp)
	}

	return fmt.Sprintf("Error with %s: %s", wr.name, wr.err.Error())
}

// JSON format for WeatherUnderground
type WeatherUnderground struct {
	CurrentObservation struct {
		Temp float64 `json:"temp_F"`
	} `json:"current_observation"`
}

func (wu WeatherUnderground) GetWeather(zipcode string, ch chan WeatherResult) {
	results := WeatherResult{name: "WeatherUndergound"}
	u, err := url.Parse("http://api.wunderground.com/api/402eb5860a9c551a/conditions/q/" + zipcode + ".json")

	if err != nil {
		results.err = err
		ch <- results
		return
	}

	err = getJSONFromHTTPCall(fmt.Sprint(u), &wu)

	if err != nil {
		results.err = err
		ch <- results
		return
	}

	results.temp = wu.CurrentObservation.Temp

	ch <- results
}

// JSON format for Aeris
type Aeris struct {
	Response []struct {
		Observation struct {
			Temp float64 `json:"tempF"`
		} `json:"ob"`
	} `json:"response"`
}

func (a Aeris) GetWeather(zipcode string, ch chan WeatherResult) {
	results := WeatherResult{name: "Aeris"}
	u, err := url.Parse("http://api.aerisapi.com/observations/closest")

	if err != nil {
		results.err = err
		ch <- results
		return
	}

	q := u.Query()
	q.Set("p", zipcode)
	q.Set("client_id", "HrJPyAxaEY9KhLDrCc6s5")
	q.Set("client_secret", "gkvSofG2W99XSO3OAaec7JM4PO4eTPN1687klztP")
	u.RawQuery = q.Encode()

	err = getJSONFromHTTPCall(fmt.Sprint(u), &a)

	if err != nil {
		results.err = err
		ch <- results
		return
	}

	if len(a.Response) == 0 {
		results.err = errors.New("Invalid response from Aeris")
		ch <- results
		return
	}

	results.temp = a.Response[0].Observation.Temp
	ch <- results
}

// JSON format for OpenWeatherMap
type OpenWeatherMap struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (owm OpenWeatherMap) GetWeather(zipcode string, ch chan WeatherResult) {
	results := WeatherResult{name: "OpenWeatherMap"}

	u, err := url.Parse("http://api.openweathermap.org/data/2.5/weather")
	if err != nil {
		results.err = err
		ch <- results
		return
	}

	q := u.Query()
	q.Set("zip", zipcode+",us")
	q.Set("units", "imperial")
	q.Set("appid", "eaa4a5db6e274fbbe0620db2196f07ad")
	u.RawQuery = q.Encode()

	err = getJSONFromHTTPCall(fmt.Sprint(u), &owm)

	if err != nil {
		results.err = err
		ch <- results
		return
	}

	results.temp = owm.Main.Temp
	ch <- results
}

// performs the HTTP GET call and returns the JSON into the passed in interface object
func getJSONFromHTTPCall(fullyFormedURL string, apiResponse interface{}) error {
	resp, err := http.Get(fullyFormedURL)

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			errorLogger.Println("Error closing http body")
		}
	}()

	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&apiResponse)
	if err != nil {
		return err
	}

	return nil
}
