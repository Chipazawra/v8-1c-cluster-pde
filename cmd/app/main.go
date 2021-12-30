package main

import (
	"log"

	"github.com/Chipazawra/v8-1c-cluster-pde/internal/app"
)

func main() {
	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
