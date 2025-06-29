package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type SimpleWeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

const (
	currentWeather = "https://api.openweathermap.org/data/2.5/weather?"
)

func main() {

	api := os.Getenv("OPEN_WEATHER_API")
	if api == "" {
		log.Fatal("Error: Open Weather API key has not been set")
	}

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	city := strings.Join(os.Args[1:], " ")

	weatherData, err := getTemp(city, api)
	if err != nil {
		log.Fatalf("Error fetching temperature for %s: %v", city, err)
	}

	fmt.Printf("Current temperature in %s: %.1f°C\n", weatherData.Name, weatherData.Main.Temp)

}

func getTemp(city string, api string) (SimpleWeatherResponse, error) {

	url := fmt.Sprintf("%sq=%s&appid=%s&units=metric", currentWeather, city, api)

	res, err := http.Get(url)
	if err != nil {
		return SimpleWeatherResponse{}, fmt.Errorf("failed to make http request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return SimpleWeatherResponse{}, fmt.Errorf("weather api returned a non-okay status. %d:%s", res.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return SimpleWeatherResponse{}, fmt.Errorf("failed to read body: %w", err)
	}

	var weatherResponse SimpleWeatherResponse
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return SimpleWeatherResponse{}, fmt.Errorf("failed to parse json: %w", err)
	}

	return weatherResponse, nil

}
