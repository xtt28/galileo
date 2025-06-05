package service

import (
	"encoding/json"
	"log"
	"net/http"
)

const ipAPIEndpoint = "http://ip-api.com/json"

var userLocation *Location

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func GetLocation() (Location, error) {
	if userLocation != nil {
		log.Println("using cached user location")
		return *userLocation, nil
	}

	res, err := http.Get(ipAPIEndpoint)
	if err != nil {
		return Location{}, err
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	var data Location

	err = dec.Decode(&data)
	if err != nil {
		return Location{}, err
	}

	userLocation = &data
	return data, nil
}
