package main

import (
	"log"

	"github.com/Grubblex/weather-api/services"
)

func main() {

	log.Println("Starting parser")
	services.Parser() 
}