package main

import (
	"fmt"

	"github.com/xtt28/galileo/app"
	"github.com/xtt28/galileo/service"
)

func main() {
	loc, err := service.GetLocation()
	if err != nil {
		fmt.Println(err)
	}

	forecast, err := service.GetCurrentWeather()
	if err != nil {
		panic(err)
	}
	temp := forecast.Properties.Timeseries[0].Data.Instant.Details.AirTemperature
	fmt.Println(*temp)

	fmt.Printf("%v\n", loc)
	app.Run()
}
