package main

import (
	"fmt"
	"os"
	"weatherapis"
)

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Println("Usage whatsthetemp <zipcode>")
		os.Exit(1)
	}

	zipcode := os.Args[1]

	apis := []weatherapis.WeatherAPI{weatherapis.WeatherUnderground{},
		weatherapis.Aeris{},
		weatherapis.OpenWeatherMap{}}

	ch := make(chan weatherapis.WeatherResult, len(apis))

	// get weather info
	for _, api := range apis {
		go api.GetWeather(zipcode, ch)
	}

	var countValidResults, totalTemp, minTemp, maxTemp float64
	needInitialValue := true

	for _ = range apis {
		result := <-ch
		fmt.Println(result)

		if result.GetError() != nil {
			continue
		}

		countValidResults++
		totalTemp += result.GetTemp()

		if needInitialValue {
			minTemp = result.GetTemp()
			maxTemp = minTemp
			needInitialValue = false

			continue
		}

		if result.GetTemp() < minTemp {
			minTemp = result.GetTemp()
		}

		if result.GetTemp() > maxTemp {
			maxTemp = result.GetTemp()
		}
	}

	fmt.Printf("\nAverage Temp: %.2f℉  ", totalTemp/countValidResults)
	fmt.Printf("Min Temp: %.2f℉  ", minTemp)
	fmt.Printf("Max Temp: %.2f℉\n", maxTemp)
}
