package weatherapis_test

import (
	"fmt"
	"weatherapis"
	"testing"
)



func TestGetWeather(t *testing.T) {
	ch := make(chan weatherapis.WeatherResult)

	api := weatherapis.MockWeatherAPI{}
	go api.GetWeather("no zipcode", ch)

	wr := <-ch

	if wr.GetError == nil {
		t.Error("wr.GetError should have value")
	}

	if wr.GetTemp() != 50 {
		t.Error("wr.GetTemp() does not equal 50")
	}

	expected := "Error with Mock Weather API: Testing API"
	actual := fmt.Sprint(wr)
	if expected != actual {
		t.Errorf("Expected '%s' but got '%s'", expected, actual)
	}
}