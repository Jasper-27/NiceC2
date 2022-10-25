package main

import (
	"log"

	"github.com/emersion/go-autostart"
)

func main() {
	app := &autostart.App{
		Name:        "test",
		DisplayName: "Just a Test App",
		Exec:        []string{"sh", "-c", "echo autostart >> ~/autostart.txt"},
	}

	if app.IsEnabled() {
		log.Println("App is already enabled, removing it...")

		if err := app.Disable(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Enabling app...")

		if err := app.Enable(); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Done!")
}
