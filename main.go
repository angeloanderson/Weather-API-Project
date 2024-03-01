package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Coordinates for Memphis, TN
	latitude := "35.1495"
	longitude := "-90.0490"

	// Make a GET request to the National Weather Service API for Memphis coordinates
	url := fmt.Sprintf("https://api.weather.gov/points/%s,%s", latitude, longitude)
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

	// Print the response body
	fmt.Println("Weather Data:", string(body))

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
