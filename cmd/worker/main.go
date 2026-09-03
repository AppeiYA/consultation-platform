package main

import (
	"log"

	app "github.com/AppeiYA/consultation-platform/internal"
)

func main() {
	w, err := app.NewWorker()
	if err != nil {
		log.Fatal(err)
	}

	if err := w.Run(); err != nil {
		log.Fatal(err)
	}
}
