package main

import (
	"fmt"
	"log"

	"github.com/denisbrodbeck/machineid"
	"github.com/emersion/go-autostart"
)

func main() {

	var TARGET_WORKING_DIRECTORY = "~/Desktop/NiceC2"

	// NodeID := generateGUID()
	NodeID, _ := machineid.ID()
	fmt.Println(NodeID)

	app := &autostart.App{
		Name:        "NiceC2",
		DisplayName: "NiceC2 command agent",
		Exec:        []string{"sh", "-c", "say 'The NiceC2 process has started' "},
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
