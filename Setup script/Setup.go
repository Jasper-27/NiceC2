package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	"github.com/emersion/go-autostart"
)

func main() {

	// var TARGET_WORKING_DIRECTORY = "~/Desktop/NiceC2"

	// NodeID := generateGUID()
	NodeID, _ := machineid.ID()
	fmt.Println(NodeID)

	// Get the users home directory
	home, _ := os.UserHomeDir()
	//testFile := filepath.Join(home, "GitHub", "NiceC2", "Agent", "agent") //Moved for cross compatability while testing

	var testFile string = ""

	// Selecting which shell to use
	var shell string
	if runtime.GOOS == "windows" {
		shell = "powershell.exe"
		testFile = filepath.Join(home, "GitHub", "NiceC2", "Agent", "agent.exe")
	} else {
		shell = "sh"
		testFile = filepath.Join(home, "GitHub", "NiceC2", "Agent", "agent")
	}

	commandString := testFile
	fmt.Println(commandString)

	app := &autostart.App{
		Name:        "NiceC2",
		DisplayName: "NiceC2 command agent",
		Exec:        []string{shell, commandString},
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
