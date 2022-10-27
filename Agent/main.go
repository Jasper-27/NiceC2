package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var command_server string = "localhost:8081"

//global stuff for shortcuts
var p = fmt.Println

func main() {

	test()

}

func test() {
	log.Println("This is a test")

	out, errorMessage := runCommand("say hello")
	if errorMessage != "" {
		log.Println(errorMessage)
		return
	} else {
		log.Println(out)
	}

}

func runCommand(command string) (outString string, errorMessage string) {

	var shell string
	errorMessage = ""

	// Selecting which shell to use
	if runtime.GOOS == "windows" {
		shell = "powershell.exe"
	} else {
		shell = "sh"
	}

	// Change directory
	if strings.HasPrefix(command, "cd ") {

		dir := command[3:] // get the first three chars

		os.Chdir(dir)

		p(dir)

		// run command, and if it causes an error create an error
		out, err := exec.Command(shell, "-c", "pwd").Output()
		if err != nil {
			p(err.Error())
			errorMessage = err.Error()

			return
		}

		outString = string(out)
		return
	}

	// run command, and if it causes an error create an error
	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {
		p(err.Error())
		errorMessage = err.Error()

		return
	}

	outString = string(out)

	return

}
