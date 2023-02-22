package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
)

// Setting the command server
// var command_server string = "http://192.168.0.69:8081"

var command_server string = "https://localhost:8081"

// var command_server string = "http://192.168.0.29:8081"

func main() {

	// This allows us to use a self signed certificate.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var task_create_request = `{ NodeID : "FCB85CB9-9452-539B-9988-48A4C5E3DFD3", Task : "run command", Details : "touch sent-from-api", Key : "2" }`

	r, err := http.NewRequest("POST", command_server+"/creat_task", bytes.NewBuffer([]byte(task_create_request)))
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
