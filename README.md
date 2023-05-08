# NiceC2
Comp3000 Final Year Project.

Trello: https://trello.com/b/cnf4AldV/nicec2 

# Description 

NiceC2 is a system for managing multiple computers on different networks. The system uses techniques often employed in command and control situations, hence the name. 

The system is designed for more technical people, who may have a collection of systems either for their own use, or for a team/group. NiceC2 will allow the user to perform updates, run scripts, and perform a remote shell on the managed devices, as long as they can talk to the command server. 

The NiceC2 server, and console are designed to be deployed on a Linux server. The agent has been designed Linux first, but is cross platform (with some missing functionality on Windows and MacOS)


# Setup 



## ServerCode 

The server is designed to be deployed on a Linux server. 

There is a setup script within the ServerCode folder, running this setup script will install the server and console, and make the server auto-start when the system boots up. 

This install script will display the server.crt. It's important to copy this out, as it is necisary for configering the agent. 


## Agent 

When deploying the NiceC2 agent the first step is to configure the code to talk to the correct command server. This is hard-coded into the agent, to avoid tampering after install. The URl for the command server can be found near the top of the main.go file. Simply replace this with the URL of the new command server. 

Next you will need to create a .crt file. This should be populated with the certificate displayed in the server setup process. 

The agent can be built with the command `go build`. Once built the executable can be coppied onto the target device, and executed on that device. The executable will copy itself and the certificate to an appropriate place in file structure, and then set itself to auto-run 

### Linux
On linux the executable will be copied to `/usr/local/bin/NiceC2/`, and a systemd service will be created to run the executable as root on system startup. The systemd service will be called NiceC2_agent.service

This service can be removed with the command `systemd disable NiceC2_agent.service`
### Windows 

On windows the executable will be copied to `C:\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\NiceC2\\"`. The program will start once the user runs in, and will be run under that user. 

### MacOS 
On MacOS the exectutable will be copied to `/Users/[current user]/Library/NiceC2/`. the program will run as the user it was installed under automatically when they log in. 

It is possible to remove the autostart by deleting the executable, or by removeing the NiceC2_agent.plist file located in `/Users/[username]/Library/LaunchAgents`. 




### Requirments 

One of the dependencies for the setup script requires GCC to be installed on the machine. This is a pain on Windows, but [Here](https://code.visualstudio.com/docs/cpp/config-mingw) is a link to some instructions.  

### Dependencies 
- github.com/denisbrodbeck/machineid
- github.com/emersion/go-autostart
- github.com/fatih/color
- github.com/gosuri/uitable

