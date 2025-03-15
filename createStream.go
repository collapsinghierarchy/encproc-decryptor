package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
)

//Example of a JWT-token for testing
//eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiMTIzNDUiLCJleHAiOjE3MzY2MDc4NDIsImlhdCI6MTczNjUyMTQ0MiwiaXNzIjoieW91ci1pc3N1ZXItbmFtZSJ9.xnVW-FmrldzR_f1fnfdKOqNtRRghKPF8IFTzcrHerVs

type he struct {
	sk     *rlwe.SecretKey
	pk     *rlwe.PublicKey
	params bgv.Parameters
}

func main() {
	// Parse JWT token from command line
	token := flag.String("token",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiMTIzNDUiLCJleHAiOjE3NDIxMjYwNDQsImlhdCI6MTc0MjAzOTY0NCwiaXNzIjoieW91ci1pc3N1ZXItbmFtZSJ9.lvP2BItbQsOa4aYxVpR6tttdptC9tCooLEmi3kaiJbg", "JWT token for API authentication")
	apiURL := flag.String("url", "http://217.154.80.44:8080/create-stream", "API URL")
	outputFile := flag.String("output", "keypair.json", "Output file to store ID, SK, and PK")
	flag.Parse()

	if *token == "" {
		log.Fatal("JWT token is required. Use -token to specify it.")
	}

	// Initialize HE struct and generate keypair
	heInstance := &he{}
	heInstance.GenerateKeypair()
	fmt.Println("Keypair generated successfully.")

	// Serialize the public key to binary
	pkBinary, err := heInstance.pk.Value.MarshalBinary()
	if err != nil {
		log.Fatal("failed to serialize public key: %v", err)
		os.Exit(1)
	}

	// Make API request
	id, err := createStream(*apiURL, *token, pkBinary)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Stream created successfully. ID: %s\n", id)

	// Store ID, SK, and PK in a file
	err = storeKeypair(*outputFile, id, heInstance.sk, heInstance.pk)
	if err != nil {
		log.Fatalf("Failed to store keypair: %v", err)
	}
	fmt.Printf("Keypair stored successfully in %s\n", *outputFile)
}

func (init *he) GenerateKeypair() {
	init.params = setupParams()

	// Key Generator
	kgen := rlwe.NewKeyGenerator(init.params)

	// Secret and Public Key
	init.sk = kgen.GenSecretKeyNew()
	init.pk = kgen.GenPublicKeyNew(init.sk)
}

func createStream(apiURL, token string, publicKey []byte) (string, error) {
	// Encode the public key as base64.
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)

	// Prepare the JSON payload with the "pk" field.
	payload := map[string]string{
		"pk": publicKeyBase64,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create a new POST request with the JSON payload.
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// Include the token as a Bearer token in the Authorization header.
	req.Header.Set("Authorization", "Bearer "+token)

	// Execute the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the API returned an OK status.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read and parse the response.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Expected response format: {"message": "Token Valid", "id": "some-id"}
	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	// Validate the response message.
	if message, ok := response["message"]; !ok || message != "Token Valid" {
		return "", errors.New("invalid token or unexpected API response")
	}

	// Retrieve the id from the response.
	id, ok := response["id"]
	if !ok {
		return "", errors.New("response did not include an ID")
	}

	return id, nil
}

// storeKeypair saves the ID, secret key, and public key as Base64-encoded binary in a JSON file.
func storeKeypair(filename, id string, sk *rlwe.SecretKey, pk *rlwe.PublicKey) error {
	// Serialize the secret key to binary
	skBinary, err := sk.Value.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to serialize secret key: %v", err)
	}

	// Serialize the public key to binary
	pkBinary, err := pk.Value.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to serialize public key: %v", err)
	}

	// Encode the binary data to Base64
	skBase64 := base64.StdEncoding.EncodeToString(skBinary)
	pkBase64 := base64.StdEncoding.EncodeToString(pkBinary)

	// Create the data structure to write to JSON
	data := map[string]interface{}{
		"id": id,
		"sk": skBase64, // Base64-encoded secret key
		"pk": pkBase64, // Base64-encoded public key
	}

	// Marshal the data to JSON format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize keypair to JSON: %v", err)
	}

	// Write the JSON data to the specified file
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write keypair to file: %v", err)
	}

	return nil
}
