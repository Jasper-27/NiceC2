package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type New_Task struct {
	NodeID  string
	Action  string
	Content string
}
type Task struct {
	TaskID   string
	NodeID   string
	Action   string
	Content  string
	Progress string
	Result   string
}
type node struct {
	ID             string
	Hostname       string
	Platform       string
	First_Check_In string
	Last_Check_In  string
}

type Task_Response_Response struct {
	TaskID string
}

var nodes []node

var tasks []Task

// Points to local host. because network security is hard.
var command_server string = "https://localhost:8081"

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to NICE C2")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("help", text) == 0 {
			fmt.Println("ls 			- List all nodes")
			fmt.Println("tasks <node> 		- Display the tasks associated with a specific device. Leave blank to show all tasks")
			fmt.Println("run <node> 		- Run a single command on a Node")
			fmt.Println("shutdown <node> 	- shutdown device")
			fmt.Println("reboot <node>		- reboot")
			fmt.Println("Exit 			- Exit the NiceC2 interface")

		}
		if strings.Compare("exit", text) == 0 {
			fmt.Println("Goodbye!")
			return
		}

		if strings.Compare("ls", text) == 0 {
			display_nodes()

		}

		if strings.Compare("tasks", text) == 0 {
			display_tasks()

		} else if strings.HasPrefix(text, "tasks ") {

			node := text[6:]
			display_task_by_node(node)

		}

		if strings.HasPrefix(text, "run") {
			node := text[4:]

			handle_run(node)

		}

		if strings.HasPrefix(text, "shutdown") {
			node := text[9:]
			shutdown(node)
		}

		if strings.HasPrefix(text, "reboot") {
			node := text[7:]
			reboot(node)
		}

	}

}

func shutdown(node string) {
	task_id := create_task_by_ID(node, "shutdown", "", "2")
	fmt.Println("Shutdown Task created (" + task_id + ")")
	time.Sleep(5 * time.Second) // Time is added to wait for command to get to / be run on node
	get_task_by_id(task_id)
}

func reboot(node string) {
	task_id := create_task_by_ID(node, "reboot", "", "2")
	fmt.Println("Reboot Task created (" + task_id + ")")
	time.Sleep(5 * time.Second) // Time is added to wait for command to get to / be run on node
	get_task_by_id(task_id)
}

func handle_run(node string) {

	// split_args := strings.Split(args, " ")

	// node := split_args[1]

	var command string

	fmt.Print("Enter command here: ")
	sub_reader := bufio.NewReader(os.Stdin)
	command, _ = sub_reader.ReadString('\n')
	// convert CRLF to LF
	command = strings.Replace(command, "\n", "", -1)

	task_id := create_task_by_ID(node, "run command", command, "2")

	fmt.Println("Waiting for command")

	// Waiting three seconds for command to complete
	time.Sleep(3 * time.Second)

	get_task_by_id(task_id)

}

func display_tasks() {

	// need to have an up to date node list
	get_nodes()

	fmt.Println("Displaying tasks")
	get_tasks()
	fmt.Println("##############################################")

	// // Displays the tasks in a sort of table thing. needs to be done better
	for _, task := range tasks {

		node_index, err := Node_index_from_node_ID(task.NodeID)
		if err != "" {
			fmt.Println(err)
		}

		// fmt.Println()
		fmt.Println("TaskID: 	" + task.TaskID)
		fmt.Println("Hostname: 	" + nodes[node_index].Hostname)
		fmt.Println("Result:")
		fmt.Println("----------")
		fmt.Println(task.Result)
		fmt.Println("##############################################")
	}
}

func Node_index_from_node_ID(nodeId string) (int, string) {
	for index, node := range nodes {
		if nodeId == node.ID {
			return index, ""
		}
	}
	return 0, "Can't find Node"
}

func NodeID_from_Hostname(input string) (string, string) {
	for _, node := range nodes {
		if input == node.Hostname {
			return node.ID, ""
		}
	}

	return "", ""
}

// Function to display the tasks assigned with the nodes. Takes either NodeID or Hostname as an argument
func display_task_by_node(NodeID string) {

	fmt.Println("Showing tasks for " + NodeID)

	get_nodes()
	get_tasks()

	// Support for both NodeID and Hostname
	var found bool = false
	for _, node := range nodes {

		if node.ID == NodeID {
			fmt.Println("Node ID detected")
			found = true
		}

		if node.Hostname == NodeID {
			fmt.Println("Hostname detected")
			NodeID = node.ID
			found = true
		}
	}
	if found == false {
		fmt.Println("No nodes found")
	}

	// output the task details
	for _, task := range tasks {
		if task.NodeID == NodeID {
			fmt.Println("######################################")
			fmt.Println("Task ID:		" + task.TaskID)
			fmt.Println("Action:			" + task.Action) // No idea why two tabs 🤷
			fmt.Println("Argument: 		" + task.Content)
			fmt.Println("Progress: 		" + task.Progress)
			fmt.Println("Result: ")
			fmt.Println("----")
			fmt.Println(task.Result)
			fmt.Println("======================================")

		}
	}
}

func display_nodes() {
	get_nodes()

	// Displays the nodes in a sort of table thing. needs to be done better
	for _, node := range nodes {
		fmt.Println("ID : ", node.ID, "	| Hostname: ", node.Hostname, "	 | Platform: ", node.Platform)
	}
}
func get_nodes() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.NewRequest("", command_server+"/get_nodes", bytes.NewBuffer([]byte("")))
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
		fmt.Println("Error sending the commands response back")
		return
	}

	// fmt.Println(res.Body)

	API_response, err := io.ReadAll(res.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	API_response_string := string(API_response)

	// Read the json string, and set the big array of nodes to it
	json.Unmarshal([]byte(API_response_string), &nodes)

}

func create_task_by_ID(nodeID string, task string, arg string, key string) string {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var task_create_request = `{
    "NodeID": "` + nodeID + `",
    "Task": "` + task + `",
    "Details": "` + arg + `",
    "Key": "2"
}`

	r, err := http.NewRequest("", command_server+"/create_task", bytes.NewBuffer([]byte(task_create_request)))
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

	// Decode the json body
	decoder := json.NewDecoder(res.Body)
	var response Task_Response_Response
	err2 := decoder.Decode(&response)
	if err2 != nil {
		panic(err2)
	}

	fmt.Println("Task submitted.  	TaskID: " + response.TaskID)

	return response.TaskID

}

func get_task_by_id(input string) {

	get_tasks()

	// // Displays the tasks in a sort of table thing. needs to be done better
	for _, task := range tasks {
		if task.TaskID == input {
			fmt.Println("######################################")
			fmt.Println("Task ID:		" + task.TaskID)
			fmt.Println("Action:			" + task.Action) // No idea why two tabs 🤷
			fmt.Println("Argument: 		" + task.Content)
			fmt.Println("Progress: 		" + task.Progress)
			fmt.Println("Result: ")
			fmt.Println("----")
			fmt.Println(task.Result)
			fmt.Println("======================================")

		}
	}

}

func get_tasks() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.NewRequest("", command_server+"/get_tasks", bytes.NewBuffer([]byte("")))
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
		fmt.Println("Error sending the commands response back")
		return
	}

	// Read the response
	API_response, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println(err2)
	}
	API_response_string := string(API_response)

	// Unmarshal the response into the tasks array
	err3 := json.Unmarshal([]byte(API_response_string), &tasks)
	if err3 != nil {
		fmt.Println(err2)
	}

}
