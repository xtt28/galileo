package main

import (
	"log"

	"github.com/xtt28/galileo/app"
	"github.com/xtt28/galileo/service"
)

func main() {
	loc, err := service.GetLocation()
	if err != nil {
		log.Println("could not get user current location")
		log.Println(err)
	}

	log.Printf("identified user current location: %v\n", loc)
	app.Run()
}
