package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type WeatherResponse struct {
	Properties struct {
		Periods []struct {
			Temperature   interface{} `json:"temperature"`
			WindSpeed     string      `json:"windSpeed"`
			ShortForecast string      `json:"shortForecast"`
		} `json:"periods"`
	} `json:"properties"`
}

func main() {
	// Make a GET request to the National Weather Service API for Memphis coordinates
	url := "https://api.weather.gov/gridpoints/MEG/42,67/forecast"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Unmarshal JSON response
	var weatherResp WeatherResponse
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Extract current weather data
	if len(weatherResp.Properties.Periods) > 0 {
		currentWeather := weatherResp.Properties.Periods[0]
		fmt.Println("Temperature:", currentWeather.Temperature)
		fmt.Println("Wind Speed:", currentWeather.WindSpeed)
		fmt.Println("Forecast:", currentWeather.ShortForecast)
	} else {
		fmt.Println("No weather data available")
		return
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	// Create a context
	ctx := context.Background()

	// Store weather data in Redis
	err = rdb.Set(ctx, "weather_data", string(body), 0).Err()
	if err != nil {
		fmt.Println("Error storing data in Redis:", err)
		return
	}

	fmt.Println("Weather data stored in Redis")
}
