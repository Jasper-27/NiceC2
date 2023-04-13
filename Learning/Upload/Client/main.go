package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func uploadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
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
	req, err := http.NewRequest("POST", "http://localhost:8083/upload", body)
	if err != nil {
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
			return fmt.Errorf("Failed to upload file. Status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("Failed to upload file. Status code: %d. Response body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func main() {
	filepath := "/Users/jasper/Desktop/uniLaptopBackground.png"

	fmt.Println("Uploading file at path:", filepath)

	err := uploadFile(filepath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("File uploaded successfully!")
}
