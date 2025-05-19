package app

import (
	"github.com/xtt28/galileo/config"

	"log"
)

func Run() {
	config.ReadConfig()
	log.Println("hello, world")
}
