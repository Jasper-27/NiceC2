package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file name from the request URL
	filename := r.URL.Path[len("/download/"):]

	// Check if the file exists
	if _, err := os.Stat("payloads/" + filename); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Open the file
	file, err := os.Open("payloads/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the response header
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", fileStat.Size()))

	// Copy the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Set the server port
	port := "8080"

	// Set the server routes
	http.HandleFunc("/download/", downloadHandler)

	// Start the server
	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
