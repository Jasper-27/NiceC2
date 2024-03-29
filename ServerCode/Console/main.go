package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"golang.org/x/term"
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

type cert struct {
	cert_content string
}

var nodes []node

var tasks []Task

// Points to local host. because network security is hard.
var command_server string = "https://localhost:8081"

func main() {

	fmt.Println()
	fmt.Println("███    ██ ██  ██████ ███████ " + color.RedString(" ██████ ") + color.RedString("██████  "))
	fmt.Println("████   ██ ██ ██      ██      " + color.RedString("██      ") + color.RedString("     ██ "))
	fmt.Println("██ ██  ██ ██ ██      █████   " + color.RedString("██      ") + color.RedString(" █████  "))
	fmt.Println("██  ██ ██ ██ ██      ██      " + color.RedString("██      ") + color.RedString("██      "))
	fmt.Println("██   ████ ██  ██████ ███████ " + color.RedString(" ██████ ") + color.RedString("███████ "))

	// fmt.Println(color.YellowString("--------------------------------------------"))
	fmt.Println("-------------------------------By Jasper Cox")

	fmt.Println("Type 'help' to see a list of commands")
	fmt.Println()

	main_loop()
}

func run_tests() {

	// test running commands

	// Testing a normal command
	test_command_1 := base64.StdEncoding.EncodeToString([]byte("touch /Users/jasper/Desktop/test_touch"))
	handle_run("Jasper's Air", test_command_1)

	// Testing command that doesn't exist
	test_command_2 := base64.StdEncoding.EncodeToString([]byte("kasdlkj akjd"))
	handle_run("Jasper's Air", test_command_2)

	// Testing command that contains special charecters
	test_command_3 := base64.StdEncoding.EncodeToString([]byte(`say "hello there"`))
	handle_run("Jasper's Air", test_command_3)

	// Test get files
	get_file("Jasper's Air", "/Users/jasper/Desktop/Not_test.txt") // test file doesn't exist
	get_file("Jasper's Air", "/Users/jasper/Desktop/test.txt")     // test file does exist
	get_file("Jasper's Air", "/Users/jasper/Desktop/test_big.zip") // test big file

	// Test sending files
	send_file("Jasper's Air", "test.txt", "/Users/jasper/Desktop/Not_real/text.txt") // broken filepath
	send_file("Jasper's Air", "test.txt", "/Users/jasper/Desktop/test.txt")          // Normal file send
	send_file("Jasper's Air", "test.txt", "/Users/jasper/Desktop/test_big.zip")      // big file

}

func main_loop() {

	var command_read bool

	// Currently selected
	var target string = ""

	reader := bufio.NewReader(os.Stdin)

	// START INPUT SECTION
	for {

		command_read = false

		fmt.Print("(" + color.BlueString(target) + ")-> ")
		text, _ := reader.ReadString('\n')

		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		// fmt.Println("|" + text + "|")

		if strings.Compare("help", text) == 0 {

			print_help_menu()

			command_read = true

		}

		if strings.Compare("exit", text) == 0 {
			// Goodbye!

			fmt.Println(color.YellowString(`
  /$$$$$$                            /$$ /$$$$$$$                      /$$
 /$$__  $$                          | $$| $$__  $$                    | $$
| $$  \__/  /$$$$$$   /$$$$$$   /$$$$$$$| $$  \ $$ /$$   /$$  /$$$$$$ | $$
| $$ /$$$$ /$$__  $$ /$$__  $$ /$$__  $$| $$$$$$$ | $$  | $$ /$$__  $$| $$
| $$|_  $$| $$  \ $$| $$  \ $$| $$  | $$| $$__  $$| $$  | $$| $$$$$$$$|__/
| $$  \ $$| $$  | $$| $$  | $$| $$  | $$| $$  \ $$| $$  | $$| $$_____/    
|  $$$$$$/|  $$$$$$/|  $$$$$$/|  $$$$$$$| $$$$$$$/|  $$$$$$$|  $$$$$$$ /$$
 \______/  \______/  \______/  \_______/|_______/  \____  $$ \_______/|__/
                                                   /$$  | $$              
                                                  |  $$$$$$/              
                                                   \______/               
`))

			os.Exit(0)
		}

		if strings.Compare("cert", text) == 0 {
			get_cert()
			command_read = true
		}

		if strings.Compare("clear", text) == 0 || strings.Compare("cls", text) == 0 {
			command_read = true

			cmd := exec.Command("clear") // create a new command to clear the terminal
			cmd.Stdout = os.Stdout       // set the command's stdout to os.Stdout
			cmd.Run()
		}

		if strings.HasPrefix(text, "use") {
			proposed_target := strings.TrimSpace(text[3:])

			target_node, err := Node_from_string(proposed_target)
			if err != nil {

				if len(proposed_target) < 1 {
					target = ""
				} else {
					fmt.Println(color.RedString("ERROR: ") + "Node not found")

				}

			} else {
				if target_node.Hostname == "" {
					target = target_node.ID
				} else {
					target = target_node.Hostname
				}

				fmt.Println("Now using node '" + target + "'")
			}

			command_read = true
		}

		if strings.Compare("ls", text) == 0 {
			display_nodes()
			command_read = true
		}

		if strings.Compare("payloads", text) == 0 {
			get_payloads_from_server()
			command_read = true
		}

		if strings.Compare("loot", text) == 0 {
			get_loot_from_server()
			command_read = true
		}

		if strings.HasPrefix(text, "tasks") {
			fmt.Println("Getting tasks")
			node := target
			if len(strings.TrimSpace(text)) > 5 {
				node = text[6:]
			}

			// Checks the node exists. If it does do the thing
			if check_node(node) != false {
				display_tasks_by_node(node)
			} else {
				fmt.Println(color.RedString("ERROR: ") + "Node not found")
			}
			command_read = true
		}

		// Running a custom command

		if strings.HasPrefix(text, "!") {
			command_read = true

			if target == "" {
				fmt.Println(color.RedString("ERROR: ") + "No node specified. Please use run [node] !command if specifying node in line")
			}

			command_string := text[1:]

			encoded_command_string := base64.StdEncoding.EncodeToString([]byte(command_string))

			handle_run(target, encoded_command_string)

		}
		if strings.HasPrefix(text, "run") {
			command_read = true
			node := target

			split := strings.SplitN(text, " !", 2) // Splits the string in 2 on the 1st instance of " !"
			if len(split) < 2 {
				fmt.Println(color.RedString("ERROR: ") + "can't parse string")

			}

			command_string := split[1]

			split2 := strings.Split(split[0], "run ")

			if len(split2) == 1 {
				// then there has been no node specified

				node = target
			} else if len(split2) == 2 {
				node = split2[1]
			}

			// Encodes command (this is to help with using "")
			encoded_command_string := base64.StdEncoding.EncodeToString([]byte(command_string))

			// Checks the node exists. If it does do the thing
			if check_node(node) != false {
				handle_run(node, encoded_command_string)
			} else {
				fmt.Println(color.RedString("ERROR: ") + "Node not found")
			}

		}

		if strings.HasPrefix(text, "shutdown") {
			command_read = true
			node := target
			if len(strings.TrimSpace(text)) > 8 {
				node = text[9:]
			}

			// Checks the node exists. If it does do the thing
			if check_node(node) != false {
				shutdown(node)
			} else {
				fmt.Println(color.RedString("ERROR: ") + "Node not found")
			}

		}

		if strings.HasPrefix(text, "reboot") {
			command_read = true
			node := target
			if len(strings.TrimSpace(text)) > 6 {
				node = text[7:]
			}

			// Checks the node exists. If it does do the thing
			if check_node(node) != false {
				reboot(node)
			} else {
				fmt.Println(color.RedString("ERROR: ") + "Node not found")
			}

		}

		if strings.HasPrefix(text, "send-file") {
			command_read = true

			processed_text := text[10:]

			var custom_target string = ""
			var file string
			var destination string

			split_1 := strings.Split(processed_text, "-f ")

			// If not target is specified
			if len(split_1) == 1 {

				//
				split_no_target := strings.Split(split_1[0], " -d ")
				if len(split_no_target) != 2 {
					fmt.Println(color.RedString("Error: ") + "Could not parse input. Missing -d desitnation")
					main_loop()
				}

				file = split_no_target[0]
				destination = split_no_target[1]

			} else if len(split_1) == 2 {
				split_2 := strings.Split(split_1[1], " -d ")
				if len(split_2) != 2 {
					fmt.Println(color.RedString("Error: ") + "Could not parse input. Missing -d desitnation")
					main_loop()
				}

				custom_target = strings.TrimSpace(split_1[0]) // trim space is needed as splitting ads a space on the end
				file = split_2[0]
				destination = split_2[1]

			} else {
				fmt.Println("Error splitting String")
				return
			}

			if custom_target != "" {

				fmt.Println("Downloading file " + file + " to " + destination + " on " + custom_target)

				// Checks the node exists. If it does do the thing
				if check_node(custom_target) == true {
					send_file(custom_target, file, destination)
				} else {
					fmt.Println(color.RedString("ERROR: ") + "Node not found")
				}

			} else {
				fmt.Println("Downloading file " + file + " to " + destination + " on " + target)
				send_file(target, file, destination)
			}

		}

		if strings.HasPrefix(text, "get-file") {
			command_read = true
			text := text[9:]

			var node string
			var path string

			split := strings.Split(text, "-p ")

			if len(split[0]) > 1 {
				node = strings.TrimSpace(split[0])
			} else {
				node = target
			}

			path = split[1]

			// checks the node is valid before sending.
			if check_node(node) != false {
				get_file(node, path)
			} else {
				fmt.Println(color.RedString("ERROR: ") + "Node not found")
			}

		}

		if command_read == false {

			fmt.Println(color.RedString("ERROR: ") + "Command not recognised")
		}

	}

}

/// END INPUT SECION

type help_message struct {
	command, description string
}

func print_help_menu() {
	table := uitable.New()
	table.MaxColWidth = 20

	x := find_terminal_width()

	if x > 10 {
		table.MaxColWidth = uint(x)/2 - 5 // Devides x by 2, and rounds down if it's odd
	}
	table.Wrap = true

	fmt.Println()
	table.AddRow("COMMAND", "DESCRIPTION")
	table.AddRow(color.RedString("-------"), color.RedString("------------"))
	table.AddRow("ls", "List all nodes")
	table.AddRow("use [node]", "Set node you are working on")
	table.AddRow("tasks [node]", "view the task queue for that node. Leaving blank will print all tasks")
	table.AddRow("run [node] ![command]", "run a single command on a node, specifying the node in line")
	table.AddRow("![command]", "run a command on a node, only works if node is already being used")
	table.AddRow("shutdown [node]", "ask a node to shutdown")
	table.AddRow("reboot [node]", "ask a node to reboot")
	table.AddRow("send-file [node] -f [filename] -d [destination file path]", "Send a file from the server to the node")
	table.AddRow("get-file [node] -p [file path on node]", "Get a file from a node, and store it on the server")
	table.AddRow("payloads", "List all the payloads available in the payloads folder")
	table.AddRow("loot", "List the files retrieved from nodes")
	table.AddRow("cert", "Show the cert needed for new nodes")
	table.AddRow("exit", "Exit the NiceC2 command line")

	fmt.Println(table)

	fmt.Println()
	fmt.Println("If a node has been sected with 'use', then it doesn't need to be specified in other commands")
	fmt.Println("Nodes in use will show up in the prompt")
	fmt.Println()
}

func find_terminal_width() int {
	if !term.IsTerminal(0) {
		return 0
	}

	width, _, err := term.GetSize(0)
	if err != nil {
		return 0
	}

	return width
}

func parse_get_file(input string) (string, string, error) {
	// Split by -f first
	parts := strings.Split(input, " -p ")
	if len(parts) != 2 {
		return "", "", errors.New("Invalid input: -p needs to come first")
	}

	return parts[0], parts[1], nil
}

func get_file(node string, path string) {

	task_id := create_task_by_ID(node, "get-file", path, "2")
	fmt.Println("get-file Task created (" + task_id + ")")
	time.Sleep(5 * time.Second) // Time is added to wait for command to get to / be run on node
	get_task_by_id(task_id)
}

func get_cert() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.NewRequest("", command_server+"/show_cert", bytes.NewBuffer([]byte("")))
	if err != nil {
		// panic(err)
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		// panic(err)
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}
	API_response, err := io.ReadAll(res.Body)
	// b, err := ioutil.ReadAll(resp.Body)  Go.1.15 and earlier
	if err != nil {
		log.Fatalln(err)
	}

	API_response_string := string(API_response)

	fmt.Println(API_response_string)

}

func send_file(node string, file string, path string) {

	// Combining the data
	data := file + " || " + path

	task_id := create_task_by_ID(node, "send-file", data, "2")
	fmt.Println("send-file Task created (" + task_id + ")")
	time.Sleep(5 * time.Second) // Time is added to wait for command to get to / be run on node
	get_task_by_id(task_id)

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

func handle_run(node string, command_string string) {

	task_id := create_task_by_ID(node, "run command", command_string, "2")

	fmt.Println("Waiting for command reply")

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
func display_tasks_by_node(NodeID string) {

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

			// Handles a base64 encoded string
			if isBase64(task.Content) == true {

				decoded, err := base64.StdEncoding.DecodeString(task.Content)
				if err != nil {
					fmt.Println("Error decoding base64:", err)
					return
				}
				fmt.Println("Argument: 		" + string(decoded))
			} else {
				fmt.Println("Argument: 		" + task.Content)
			}

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

	table := uitable.New()
	table.MaxColWidth = 20

	x := find_terminal_width()

	if x > 10 {
		table.MaxColWidth = uint(x)/3 - 5 // Devides x by 2, and rounds down if it's odd
	}
	table.Wrap = true

	table.AddRow(color.GreenString("ID"), color.GreenString("HOSTNAME"), color.GreenString("PLATFORM"), color.GreenString("LAST CHECK IN"))
	table.AddRow(color.WhiteString("--"), color.WhiteString("--------"), color.WhiteString("---------"), color.WhiteString("-------------------"))
	// Displays the nodes in a sort of table thing. needs to be done better
	for _, node := range nodes {
		table.AddRow(node.ID, node.Hostname, node.Platform, convertToPretyyTime(node.Last_Check_In))
	}

	fmt.Println(table)

	fmt.Println("")
}

func convertToPretyyTime(datetimeStr string) string {

	processed_text := datetimeStr[:19]

	return (processed_text)
}

func check_node(input string) bool {
	_, err := Node_from_string(input)

	if err != nil {
		return false
	}

	return true

}

// Function to check if node exists
func Node_from_string(input string) (node, error) {

	// Makes sure we have the most up to date list of nodes
	get_nodes()

	// Support for both NodeID and Hostname
	var found bool = false
	var empty_node node // create an empty node for if none are found.
	for _, node := range nodes {

		if node.ID == input {
			found = true
			return node, nil
		}

		if node.Hostname == input {
			found = true
			return node, nil
		}
	}
	if found == false {
		return empty_node, errors.New("Node doesn't exist")
	}

	// This line will never be executed, but Go requires it
	return empty_node, nil
}

func get_nodes() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.NewRequest("", command_server+"/get_nodes", bytes.NewBuffer([]byte("")))
	if err != nil {
		// panic(err)
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
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

	if nodeID == "" {
		fmt.Println(color.RedString("ERROR: ") + "No node specified")
		main_loop()
		return ""
	}

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
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		// panic(err)
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
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

			// Handles a base64 encoded string
			if isBase64(task.Content) == true {

				decoded, err := base64.StdEncoding.DecodeString(task.Content)
				if err != nil {
					fmt.Println("Error decoding base64:", err)
					return
				}
				fmt.Println("Argument: 		" + string(decoded))
			} else {
				fmt.Println("Argument: 		" + task.Content)
			}

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
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
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

func get_payloads_from_server() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.NewRequest("", command_server+"/list_payloads", bytes.NewBuffer([]byte("")))
	if err != nil {
		// panic(err)
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
		return
	}

	// Read the response
	API_response, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println(err2)
	}
	API_response_string := string(API_response)

	// Converting the string into a slice of filenames.
	var slice []string
	err3 := json.Unmarshal([]byte(API_response_string), &slice)
	if err3 != nil {
		panic(err)
	}

	// Making a nice output

	fmt.Println()
	fmt.Println("Payloads stored in /usr/local/bin/NiceC2_server/payloads/")
	fmt.Println("#############################")
	fmt.Println()
	for _, item := range slice {
		fmt.Println(item)
	}

	fmt.Println()

}

func get_loot_from_server() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.NewRequest("", command_server+"/list_loot", bytes.NewBuffer([]byte("")))
	if err != nil {
		// panic(err)
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(color.RedString("ERROR: ") + "can't communicate with server.")
		main_loop()
		return
	}

	// Read the response
	API_response, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println(err2)
	}
	API_response_string := string(API_response)

	// Converting the string into a slice of filenames.
	var slice []string
	err3 := json.Unmarshal([]byte(API_response_string), &slice)
	if err3 != nil {
		panic(err)
	}

	// Making a nice output

	fmt.Println()
	fmt.Println("Loot is stored in /usr/local/bin/NiceC2_server/loot/")
	fmt.Println("#############################")
	fmt.Println()
	for _, item := range slice {
		fmt.Println(item)
	}

	fmt.Println()

}

// Checks if a string is base64 encoded
func isBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
