package main

import (
	_ "keybook/backend/migrations"
	"log"

	"github.com/pocketbase/pocketbase"
	"go.uber.org/dig"
)

func startBackend() {
	container := dig.New()
	container.Provide(pocketbase.New)

	invokeErr := container.Invoke(func(app *pocketbase.PocketBase) {
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
	})
	if invokeErr != nil {
		log.Fatal(invokeErr)
	}
}

func main() {
	startBackend()
}
