package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadFile(serverURL string, filename string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Make the request
	resp, err := http.Get(serverURL + "/download/" + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		os.Remove(filepath) // if the file isn't found. remove the empty destination
		return fmt.Errorf("failed to download file: %s", resp.Status)
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

	return nil
}

func main() {
	// Set the server URL
	serverURL := "http://localhost:8080"
	// Set the file name to download
	filename := "tests.iso"

	// Set the file path for the downloaded file
	// filePath := "example.txt"

	// Download the file
	err := downloadFile(serverURL, filename, "Downloads//"+filename)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}

	fmt.Println("File downloaded successfully!")
}
