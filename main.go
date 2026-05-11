package main

import (
	"log"

	"oojsite/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
