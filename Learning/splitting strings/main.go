package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {

	input := "del23 -d /etc/programs/somewhere -f movie.txt"

	DeviceName, file, destination, err := parse_download(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Device: " + DeviceName)
	fmt.Println("File: " + file)
	fmt.Println("Destination: " + destination)

}

func parse_download(input string) (string, string, string, error) {

	// Split by -f first
	parts1 := strings.Split(input, " -f ")
	if len(parts1) != 2 {
		return "", "", "", errors.New("Invalid input: -f needs to come first")
	}

	// Split the second part by -d
	parts2 := strings.Split(parts1[1], " -d ")
	if len(parts2) != 2 {
		return "", "", "", errors.New("Invalid input: missing -d flag")
	}

	// deviceName / file / destination
	return string(parts1[0]), string(parts2[0]), string(parts2[1]), nil

}
