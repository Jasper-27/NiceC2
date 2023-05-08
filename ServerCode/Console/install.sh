#!/bin/bash

echo "Starting Installing the console" 


# Set the name of your Golang program
PROGRAM_NAME="NiceC2"

# Set the path to the directory where you want to install your program
INSTALL_DIR="/usr/local/bin/"


# Build the Golang program
go build -o $PROGRAM_NAME 

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Error: Failed to build $PROGRAM_NAME"
  exit 1
fi

# Install the program
sudo cp $PROGRAM_NAME $INSTALL_DIR/$PROGRAM_NAME
sudo chmod +x $INSTALL_DIR/$PROGRAM_NAME


echo "Done installing the console"