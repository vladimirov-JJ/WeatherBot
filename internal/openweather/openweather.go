package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
	}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := http.Get(fmt.Sprintf(url, city, o.apiKey))
	if err != nil {
		return Coordinates{}, fmt.Errorf("error get coordinates: %w", err)
	}

	if resp.StatusCode != 200 {
		return Coordinates{}, fmt.Errorf("error code get coordinates: %d", resp.StatusCode)
	}

	var coordinatesResponse []CoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		return Coordinates{}, fmt.Errorf("error unmarshal response: %w", err)
	}

	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("error empty coordinates")
	}

	return Coordinates{
		Name: coordinatesResponse[0].Name,
		Lat:  coordinatesResponse[0].Lat,
		Lon:  coordinatesResponse[0].Lon,
	}, nil
}

func (o OpenWeatherClient) Weather(lat, lon float64) (WeatherResponse, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"
	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.apiKey))
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("error get weather: %w", err)
	}

	if resp.StatusCode != 200 {
		return WeatherResponse{}, fmt.Errorf("error code get weather: %d", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("error unmarshal weather response: %w", err)
	}

	return weatherResponse, nil
}
