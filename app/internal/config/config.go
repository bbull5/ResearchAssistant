package config

import (
	"log"
	"os"
)


var Port string

func LoadConfig() {
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	log.Println("Using port:", Port)
}