#!/bin/bash


# Get the hostname for the cirtificates 
read -p "Enter common name for certificates: " hostname

echo $hostname

# Set the name of your Golang program
PROGRAM_NAME="NiceC2_Console"

# Set the path to the directory where you want to install your program
INSTALL_DIR="/usr/local/bin/NiceC2_server"
mkdir $INSTALL_DIR


# Set the version number of your program
VERSION="1.0.0"

# Build the Golang program
go build -o $PROGRAM_NAME 

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Error: Failed to build $PROGRAM_NAME"
  exit 1
fi

# Install the program
sudo cp $PROGRAM_NAME $INSTALL_DIR/$PROGRAM_NAME
sudo chmod +x $INSTALL_DIR/*

