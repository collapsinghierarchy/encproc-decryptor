package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
)

func mainxy() {
	_, sk, pk, err := loadKeypair("keypair.json")
	if err != nil {
		panic(err)
	}
	he := &he{sk, pk, setupParams()}

	input_ar := []int{2, 2, 4, 5, 6, 7}
	ct, err := he.encrypt_input(input_ar, pk)
	if err != nil {
		panic(err)
	}
	pt, err := he.decrypt_result(sk, ct)
	if err != nil {
		panic(err)
	}
	fmt.Println(pt)
}

func (enc *he) encrypt_input(input_ar []int, pk *rlwe.PublicKey) (*rlwe.Ciphertext, error) {
	var err error

	//----------------Init PT ----------------
	pt := bgv.NewPlaintext(enc.params, enc.params.MaxLevel())
	// Encoder
	ecd := bgv.NewEncoder(enc.params)
	input := make([]uint64, enc.params.MaxSlots())

	// Set the positions corresponding to the provided indices to 1
	for index, v := range input_ar {
		input[index] = uint64(v)
	}

	// Encryptor
	encryptor := bgv.NewEncryptor(enc.params, pk)

	//-----------------Encode & Encrypt-------
	// Encodes the vector of plaintext values
	if err := ecd.Encode(input, pt); err != nil {
		panic(err)
	}
	// Encrypts the vector of plaintext values
	ct := bgv.NewCiphertext(enc.params, 1, enc.params.MaxLevel())
	err = encryptor.Encrypt(pt, ct)
	if err != nil {
		panic(err)
	}

	return ct, nil
}

func (enc *he) decrypt_result(sk *rlwe.SecretKey, ct *rlwe.Ciphertext) ([]uint64, error) {
	var err error
	// Encoder
	ecd := bgv.NewEncoder(enc.params)
	// Decryptor
	dec := bgv.NewDecryptor(enc.params, sk)
	// Decrypts the vector of plaintext values

	pt := rlwe.NewPlaintext(enc.params, enc.params.MaxLevel())
	dec.Decrypt(ct, pt)
	// Vector of plaintext values
	values := make([]uint64, enc.params.MaxSlots())
	// Encoder
	if err = ecd.Decode(pt, values); err != nil {
		panic(err)
	}
	return values, nil
}

// loadKeypair loads the ID, secret key, and public key from a JSON file where they were stored
// as Base64-encoded strings.
func loadKeypair(filename string) (string, *rlwe.SecretKey, *rlwe.PublicKey, error) {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to read keypair file: %v", err)
	}

	// Parse the JSON data
	data := make(map[string]string)
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to parse keypair JSON: %v", err)
	}

	// Extract the ID
	id, ok := data["id"]
	if !ok {
		return "", nil, nil, fmt.Errorf("missing ID in keypair file")
	}

	// Extract and decode the secret key
	skBase64, ok := data["sk"]
	if !ok {
		return "", nil, nil, fmt.Errorf("missing secret key in keypair file")
	}
	skBinary, err := base64.StdEncoding.DecodeString(skBase64)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to decode secret key: %v", err)
	}

	// Extract and decode the public key
	pkBase64, ok := data["pk"]
	if !ok {
		return "", nil, nil, fmt.Errorf("missing public key in keypair file")
	}

	pkBinary, err := base64.StdEncoding.DecodeString(pkBase64)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	// Reconstruct the secret key
	sk := new(rlwe.SecretKey)
	err = sk.Value.UnmarshalBinary(skBinary)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to reconstruct secret key: %v", err)
	}

	// Reconstruct the public key
	pk := new(rlwe.PublicKey)
	err = pk.Value.UnmarshalBinary(pkBinary)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to reconstruct public key: %v", err)
	}

	return id, sk, pk, nil
}
