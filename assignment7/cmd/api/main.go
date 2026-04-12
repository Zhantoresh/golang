package main

import (
	"assignment7/config"
	"assignment7/internal/app"
)

func main() {
	
	config.LoadConfig()

	
	app.Run()
}