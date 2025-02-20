# encproc-decryptor

**encproc-decryptor** is the client-side decryption utility designed to work in tandem with the [encproc](https://github.com/collapsinghierarchy/encproc) engine—your Encrypted Processing as a Service solution. This repository provides a set of tools and examples to help you create encrypted data streams, contribute encrypted data, and retrieve & decrypt aggregated results from the encproc engine.

> **Note:** You must have the encproc engine up and running (see [encproc](https://github.com/collapsinghierarchy/encproc)) and a valid JWT token for authentication to start experimenting on your own. 

---

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Creating a New Stream](#creating-a-new-stream)
  - [Contributing Data to a Stream](#contributing-data-to-a-stream)
  - [Retrieving Aggregated Results](#retrieving-aggregated-results)
- [How It Works](#how-it-works)
- [Troubleshooting](#troubleshooting)
- [License](#license)
- [Contributing](#contributing)

---

## Overview

The **encproc-decryptor** repository serves as a client-side complement to the encproc engine. It includes:
- A minimal Go-based API client for creating streams. (You only need a valid JWT-token)
- HTML/JavaScript examples for interacting with the engine’s endpoints.
- Sample usage instructions that demonstrate how to:
  - Create a new data stream.
  - Contribute encrypted data to an existing stream.
  - Retrieve aggregated (encrypted) results.
  - Decrypt the aggregated results on the client side.

This utility abstracts the complexities of interacting with homomorphic encryption and provides simple examples to help you integrate decryption into your workflow.

---

## Prerequisites

Before using **encproc-decryptor**, ensure that you have:
- A running instance of the [encproc](https://github.com/collapsinghierarchy/encproc) engine.
- A valid JWT token for authentication with the encproc engine.
- [Go](https://golang.org/) installed if you wish to compile or modify the Go code.
- Basic knowledge of using `curl` or a similar HTTP client for API requests.
- (Optional) A modern web browser if you plan to use the provided HTML/JavaScript examples.

---

## Installation

Clone the repository using Git:

```bash
git clone https://github.com/collapsinghierarchy/encproc-decryptor.git
cd encproc-decryptor
```

If you need to run or modify the Go code (e.g., `createStream.go`, `params.go`), make sure your Go environment is set up. Then install the dependencies:

```bash
go mod tidy
```

For client-side usage (HTML/JavaScript), simply open the HTML files in your browser or serve them via your preferred web server.

---

## Usage

### Creating a New Stream

To initialize a new data stream in the encproc engine, use the following API call. Replace `<your_jwt_token>` with your JWT token and `<your_pk>` with your public key (used by the encryption module):

```bash
curl -X POST http://localhost:8080/create-stream      -H "Authorization: Bearer <your_jwt_token>"      -d '{"pk": "<your_pk>"}'
```

This endpoint will create a new stream and return a unique stream ID that is used in subsequent operations.

---

### Contributing Data to a Stream

Once a stream is created, you can contribute encrypted data to it. Use the stream ID returned from the stream creation step and replace `<your_encrypted_data>` with the ciphertext produced by your encryption module:

```bash
curl -X POST http://localhost:8080/streams/{streamID}/contribute      -H "Authorization: Bearer <your_jwt_token>"      -d '{"encrypted_data": "<your_encrypted_data>"}'
```

This call forwards your encrypted data to the encproc engine where it will be aggregated.

---

### Retrieving Aggregated Results

After enough contributions have been made, you can query the aggregated (encrypted) results. Use the same stream ID as before:

```bash
curl -X GET http://localhost:8080/streams/{streamID}/aggregate      -H "Authorization: Bearer <your_jwt_token>"
```

The response will include:
- `ct_aggr_byte_base64`: The aggregated ciphertext (in Base64).
- `sample_size`: The number of contributions that were aggregated.

To decrypt these results on the client side, use the decryption functions provided in this repository (via the HTML/JavaScript demo or your own integration). See the [How It Works](#how-it-works) section below for more details.

---

## How It Works

1. **Stream Creation:**  
   The client (encproc-decryptor) sends a request to the encproc engine to create a new data stream. This stream is identified by a unique stream ID and associated with your public key.

2. **Data Contribution:**  
   Clients encrypt their data locally using the provided encryption libraries (from the encproc engine) and send the encrypted data to the stream endpoint.

3. **Aggregation:**  
   The encproc engine aggregates the encrypted contributions using homomorphic operations. The aggregated data is stored as an encrypted result.

4. **Decryption:**  
   On the client side, **encproc-decryptor** provides tools (e.g., HTML/JavaScript demo, Go functions) to decrypt the aggregated result. The decryption process uses the secret key (which must be kept secure and loaded on the client side) to convert the aggregated ciphertext into a readable result.

---

## Troubleshooting

- **JWT Token Issues:**  
  Ensure that your JWT token is valid and has the required permissions to create streams and contribute data.

- **Engine Connectivity:**  
  Verify that the encproc engine is running on the specified host and port (e.g., `http://localhost:8080`).

- **Encryption/Decryption Errors:**  
  If decryption fails, confirm that you are using the correct secret key and that the aggregated ciphertext is not corrupted.

- **Dependency Issues:**  
  For Go-related issues, run `go mod tidy` to ensure all dependencies are correctly installed.

---

## License

This repository is licensed under the [Apache-2.0 license](LICENSE).

---

## Contributing

Contributions to **encproc-decryptor** are welcome! Feel free to submit issues or pull requests with improvements or bug fixes. Please ensure that your contributions align with the overall design of the encproc engine and its client-side decryption process.

For major changes, please open an issue first to discuss what you would like to change.

---

## Additional Resources

- [encproc Repository](https://github.com/collapsinghierarchy/encproc) – The server-side engine that powers the encrypted processing as a service.
- [Encryption and Decryption Demo](#) – (Link to any demo or additional documentation if available.)

Happy decrypting!
