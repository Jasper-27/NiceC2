package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
)

var command_server string = "https://localhost:8081"

func main() {

	fmt.Println("client")

	// Read the server certificate from a file
	certPEM, err := ioutil.ReadFile("certificate.pem")
	if err != nil {
		panic(err)
	}

	// Parse the certificate
	serverCert, err := x509.ParseCertificate(certPEM)
	if err != nil {
		panic(err)
	}

	// Calculate the SHA-256 hash of the certificate
	certHash := sha256.Sum256(serverCert.Raw)

	// Convert the hash to a hexadecimal string
	expectedCertHash := hex.EncodeToString(certHash[:])

	// Create a custom certificate verification function that checks the certificate's hash
	verifyCert := func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) == 0 {
			return fmt.Errorf("no certificate provided by the server")
		}

		_, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return fmt.Errorf("failed to parse certificate: %v", err)
		}

		// Calculate the SHA-256 hash of the certificate
		hash := sha256.Sum256(rawCerts[0])
		fmt.Println(hash)

		// Compare the calculated hash with the expected hash
		if !bytes.Equal(hash[:], certHash[:]) {
			return fmt.Errorf("certificate hash mismatch, expected %s, got %s", expectedCertHash, hex.EncodeToString(hash[:]))
		}

		// If the certificate hash matches, return no error
		return nil
	}

	// Create a TLS configuration with our custom certificate verification function
	tlsConfig := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyCert,
	}

	// Getting the info
	hostname, _ := os.Hostname()
	platform := runtime.GOOS

	data := map[string]string{"ID": "2", "Hostname": hostname, "Platform": platform}
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error Marshalling the JSON sent by the server")
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", command_server+"/checkin", bytes.NewBuffer(json_data))
	if err != nil {
		fmt.Println("Error posting to command server")
		return
	}

	// Add the header to say that it's json
	req.Header.Add("Content-Type", "application/json")

	// Create an HTTP client with the custom transport
	// Create an HTTP client that uses our custom TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	// Bit of UI for you
	if res.StatusCode == 200 {
		// fmt.Println("Response was OK")
	} else {
		fmt.Println("Response was Not OK")
	}

	fmt.Println(res.StatusCode)
}
