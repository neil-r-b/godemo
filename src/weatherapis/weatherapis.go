package weatherapis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	} else {
		return fmt.Sprintf("Error with %s: %s", wr.name, wr.err.Error())
	}
}

// JSON format for WeatherUnderground
type WeatherUnderground struct {
	CurrentObservation struct {
		Temp float64 `json:"temp_F"`
	} `json:"current_observation"`
}

func (wu WeatherUnderground) GetWeather(zipcode string, ch chan WeatherResult) {
	url := "http://api.wunderground.com/api/402eb5860a9c551a/conditions/q/" + zipcode
	url += ".json"

	results := WeatherResult{name: "WeatherUndergound"}
	err := getJSONFromHTTPCall(url, &wu)

	if err != nil {
		results.err = err
		ch <- results
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
	url := "http://api.aerisapi.com/observations/closest?p=" + zipcode
	url += "&client_id=HrJPyAxaEY9KhLDrCc6s5&client_secret=gkvSofG2W99XSO3OAaec7JM4PO4eTPN1687klztP"

	results := WeatherResult{name: "Aeris"}
	err := getJSONFromHTTPCall(url, &a)

	if err != nil {
		results.err = err
		ch <- results
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
	url := "http://api.openweathermap.org/data/2.5/weather?zip=" + zipcode
	url += ",us&units=imperial&appid=eaa4a5db6e274fbbe0620db2196f07ad"

	results := WeatherResult{name: "OpenWeatherMap"}
	err := getJSONFromHTTPCall(url, &owm)

	if err != nil {
		results.err = err
		ch <- results
	}

	results.temp = owm.Main.Temp
	ch <- results
}

// below used for testing, but does this make sense?
type MockWeatherAPI struct {
}

func (m MockWeatherAPI) GetWeather(zipcode string, ch chan WeatherResult) {
	ch <- WeatherResult{name: "Mock Weather API", temp: 50, err: errors.New("Testing API")}
}

// performs the HTTP GET call and returns the JSON into the passed in interface object
func getJSONFromHTTPCall(fullyFormedURL string, apiResponse interface{}) error {
	resp, err := http.Get(fullyFormedURL)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		return err
	}

	return nil
}
