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

// FetchWeather fetches weather data from the National Weather Service API.
func FetchWeather(url string) ([]byte, error) {
    response, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}

// ExtractCurrentWeather extracts current weather data from the JSON response.
func ExtractCurrentWeather(body []byte) (WeatherResponse, error) {
    var weatherResp WeatherResponse
    err := json.Unmarshal(body, &weatherResp)
    if err != nil {
        return WeatherResponse{}, err
    }

    if len(weatherResp.Properties.Periods) > 0 {
        return weatherResp, nil
    }
    return WeatherResponse{}, fmt.Errorf("no weather data available")
}

// StoreWeatherData stores weather data in Redis.
func StoreWeatherData(ctx context.Context, rdb *redis.Client, key string, data []byte) error {
    err := rdb.Set(ctx, key, string(data), 0).Err()
    if err != nil {
        return err
    }
    return nil
}

func main() {
    // Make a GET request to the National Weather Service API for Memphis coordinates
    url := "https://api.weather.gov/gridpoints/MEG/42,67/forecast"
    body, err := FetchWeather(url)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Extract current weather data
    weatherResp, err := ExtractCurrentWeather(body)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Print current weather data
    currentWeather := weatherResp.Properties.Periods[0]
    fmt.Println("Temperature:", currentWeather.Temperature)
    fmt.Println("Wind Speed:", currentWeather.WindSpeed)
    fmt.Println("Forecast:", currentWeather.ShortForecast)

    // Connect to Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Redis server address
    })

    // Create a context
    ctx := context.Background()

    // Store weather data in Redis
    err = StoreWeatherData(ctx, rdb, "weather_data", body)
    if err != nil {
        fmt.Println("Error storing data in Redis:", err)
        return
    }

    fmt.Println("Weather data stored in Redis")
}
