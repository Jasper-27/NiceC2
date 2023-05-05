#!/bin/bash

# Generate self-signed CA certificate
openssl req -x509 -newkey rsa:2048 -keyout ca.key -out ca.crt -days 365 -nodes -subj '/CN=root-27.duckdns.org'

# Generate server key and CSR
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj '/CN=root-27.duckdns.org'

# Sign server CSR with CA to create server certificate
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365

# Generate client key and CSR
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj '/CN=Client'

# Sign client CSR with CA to create client certificate
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 365

# Clean up CSR files
rm *.csr

echo "Certificates generated successfully!"