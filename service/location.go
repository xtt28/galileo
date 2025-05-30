package service

import (
	"encoding/json"
	"net/http"
)

const ipAPIEndpoint = "https://ipapi.co/json"

var userLocation *Location

type Location struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func GetLocation() (Location, error) {
	if userLocation != nil {
		return *userLocation, nil
	}

	res, err := http.Get(ipAPIEndpoint)
	if err != nil {
		return Location{}, err
	}

	dec := json.NewDecoder(res.Body)
	var data Location

	err = dec.Decode(&data)
	if err != nil {
		return Location{}, err
	}

	return data, nil
}
