package service

import (
	"time"

	"github.com/zapling/yr.no-golang-client/client"
	"github.com/zapling/yr.no-golang-client/locationforecast"
)

type forecastCache struct {
	data *locationforecast.GeoJson
	expiry time.Time
}

var forecastCached *forecastCache

func GetCurrentWeather(cl *client.YrClient) (*locationforecast.GeoJson, error) {
	if forecastCached != nil && forecastCached.expiry.Before(time.Now()) {
		return forecastCached.data, nil
	}
	
	location, err := GetLocation()
	if err != nil {
		return nil, err
	}
	forecast, _, err := locationforecast.GetCompact(cl, location.Latitude, location.Longitude)
	forecastCached = &forecastCache{forecast, time.Now().Add(10 * time.Minute)}
	return forecast, err
}
