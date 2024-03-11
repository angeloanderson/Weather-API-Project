package main

import (
    "context"
    "testing"

    "github.com/go-redis/redis/v8"
)

func TestFetchWeather(t *testing.T) {
    url := "https://api.weather.gov/gridpoints/MEG/42,67/forecast"
    _, err := FetchWeather(url)
    if err != nil {
        t.Errorf("FetchWeather() returned an error: %v", err)
    }
}

func TestExtractCurrentWeather(t *testing.T) {
    sampleData := []byte(`{"properties":{"periods":[{"temperature":72,"windSpeed":"5 mph","shortForecast":"Sunny"}]}}`)
    _, err := ExtractCurrentWeather(sampleData)
    if err != nil {
        t.Errorf("ExtractCurrentWeather() returned an error: %v", err)
    }
}

func TestStoreWeatherData(t *testing.T) {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    ctx := context.Background()
    err := StoreWeatherData(ctx, rdb, "test_key", []byte("test_data"))
    if err != nil {
        t.Errorf("StoreWeatherData() returned an error: %v", err)
    }
}
