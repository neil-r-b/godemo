package weatherapis

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockWeatherAPI struct {
}

func (m MockWeatherAPI) GetWeather(zipcode string, ch chan WeatherResult) {
	ch <- WeatherResult{name: "Mock Weather API", temp: 50, err: errors.New("Testing API")}
}

func TestMockGetWeather(t *testing.T) {
	ch := make(chan WeatherResult)

	api := MockWeatherAPI{}
	go api.GetWeather("no zipcode", ch)

	wr := <-ch

	if wr.GetError == nil {
		t.Error("wr.GetError should have value")
	}

	if wr.GetError().Error() != "Testing API" {
		t.Errorf("Expected 'Testing API' but got %s", wr.GetError().Error())
	}

	if wr.GetTemp() != 50 {
		t.Errorf("Expected Temp = '50' but got %s", wr.GetTemp())
	}

	expected := "Error with Mock Weather API: Testing API"
	actual := fmt.Sprint(wr)
	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}

type MockHTTPWeatherAPI struct {
	Main struct {
		Observation struct {
			Temp float64 `json:"temp"`
		} `json:"observation"`
	} `json:"main"`

	url string
}

func (m MockHTTPWeatherAPI) GetWeather(zipcode string, ch chan WeatherResult) {
	wr := WeatherResult{name: "MockHTTPWeatherAPI"}

	err := getJSONFromHTTPCall(m.url, &m)
	if err != nil {
		wr.err = err
		ch <- wr
		return
	}

	wr.temp = m.Main.Observation.Temp
	ch <- wr
}

func TestMockHTTPGetWeather(t *testing.T) {
	ch := make(chan WeatherResult)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{ "main": { "observation": { "temp": 50.3 } } }`)
	}))

	api := MockHTTPWeatherAPI{url: ts.URL}
	go api.GetWeather("zipcode", ch)

	wr := <-ch

	if wr.name != "MockHTTPWeatherAPI" {
		t.Errorf("Expected 'MockHTTPWeatherAPI' but got %s", wr.name)
	}

	if wr.temp != 50.3 {
		t.Errorf("Expected 50.3 but got %f", wr.temp)
	}
}
