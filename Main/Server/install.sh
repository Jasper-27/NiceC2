#!/bin/bash


# Get the hostname for the cirtificates 
read -p "Enter common name for certificates: " hostname

echo $hostname

# Set the name of your Golang program
PROGRAM_NAME="NiceC2_server"

# Set the path to the directory where you want to install your program
INSTALL_DIR="/usr/local/bin/NiceC2_server/"
mkdir $INSTALL_DIR

# Making the uploads and payloads folder 

mkdir $INSTALL_DIR/payloads/
mkdir $INSTALL_DIR/uploads/

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


# Setup server certificates 
openssl req -new -subj "/C=GB/ST=Devon/CN=$hostname" -newkey rsa:2048 -nodes -keyout $INSTALL_DIR/server.key -out $INSTALL_DIR/server.csr   
openssl x509 -req -days 365 -in $INSTALL_DIR/server.csr -signkey $INSTALL_DIR/server.key -out $INSTALL_DIR/server.crt

echo "Key for clients: $INSTALL_DIR"
cat $INSTALL_DIR/server.crt

# Create a systemd service for the program
sudo tee /etc/systemd/system/$PROGRAM_NAME.service > /dev/null << EOF
[Unit]
Description=$PROGRAM_NAME

[Service]
Type=simple
ExecStart=$INSTALL_DIR/$PROGRAM_NAME

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd to detect the new service
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable $PROGRAM_NAME

# Start the service
sudo systemctl start $PROGRAM_NAME

echo "Success: Installed $PROGRAM_NAME to $INSTALL_DIR and configured it to start on boot"
