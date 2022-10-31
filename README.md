# NiceC2
Comp3000 Final Year Project.


## Links 

Trello: https://trello.com/b/cnf4AldV/nicec2 

## Description 

NiceC2 is a system for managing multiple computers on different networks. The system uses techniques often employed in command and control situations, hence the name. 

The system is designed for more technical people, who may have a collection of systems either for their own use, or for a team/group. NiceC2 will allow the user to perform updates, run scripts, and perform a remote shell on the managed devices, as long as they can talk to the command server. 


## Wait, isn’t this just malware? 


Nope. There are many situations where someone would legitimately benefit from the ability to control multiple machines in this way. For example someone who wants to keep their home lab servers updated, or a small business that can’t justify the use of a truly enterprise solution like Intune. 


## So it’s like BETEC Intune? 

Kinda, but with like Linux support, and you control the command server. 


## Notes 

### Requirments 

One of the dependencies for the setup script requires GCC to be installed on the machine. This is a pain on Windows, but [Here](https://code.visualstudio.com/docs/cpp/config-mingw) is a link to some instructions.  


### Trouble with AutoStarts 

At the moment the startup script needs to be pointed at the compiled version of the agent. This is setup just for Windows at the moment. 

