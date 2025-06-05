package service

import (
	"log"
	"net/http"
	"time"

	"github.com/zapling/yr.no-golang-client/client"
	"github.com/zapling/yr.no-golang-client/locationforecast"
)

var yrNoClient = client.NewYrClient(http.DefaultClient, "io.github.xtt28.galileo")

type forecastCache struct {
	data   *locationforecast.GeoJson
	expiry time.Time
}

var forecastCached *forecastCache

func GetCurrentWeather() (*locationforecast.GeoJson, error) {
	if forecastCached != nil && forecastCached.expiry.Before(time.Now()) {
		log.Println("fetching weather data - no living cached forecast")
		return forecastCached.data, nil
	}

	location, err := GetLocation()
	if err != nil {
		return nil, err
	}
	forecast, _, err := locationforecast.GetCompact(yrNoClient, location.Latitude, location.Longitude)
	forecastCached = &forecastCache{forecast, time.Now().Add(30 * time.Minute)}
	log.Printf("caching weather forecast until %s\n", forecastCached.expiry)
	
	return forecast, err
}
