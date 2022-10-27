package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// Sending a nice response
	fmt.Fprint(w, "Node has checked in!")

	// Decode the json body
	decoder := json.NewDecoder(req.Body)
	var node node_check_in
	err := decoder.Decode(&node)
	if err != nil {
		panic(err)
	}

	log.Println(node.ID + " Checked in")

}

func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)

	// myRouter.HandleFunc("/test", shrug).Methods("GET")
	myRouter.HandleFunc("/checkin", nodeCheckIn).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
}

func shrug() {
	fmt.Println("🤷")
}
