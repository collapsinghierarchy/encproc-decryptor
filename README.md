# encproc-decryptor

**encproc-decryptor** is the client-side decryption utility designed to work in tandem with the [encproc](https://github.com/collapsinghierarchy/encproc) engine—your Encrypted Processing as a Service solution. This repository provides a set of tools and examples to help you create encrypted data streams, contribute encrypted data, and retrieve & decrypt aggregated results from the encproc engine.

> **Note:** You must have the encproc engine up and running (see [encproc](https://github.com/collapsinghierarchy/encproc)) and a valid JWT token for authentication to start experimenting on your own. Reach to me in such a case. Currently, there is also no instance of the encproc engine running, but i plan to launch one in the near future.

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
- A minimal Go-based API client for creating streams. 
- HTML/JavaScript examples for interacting with the engine’s endpoints. You can customize them in any way you need. They are the primary target for experimentation.
- Sample usage instructions that demonstrate how to:
  - Create a new data stream.
  - Contribute encrypted data to an existing stream.
  - Retrieve aggregated (encrypted) results.
  - Decrypt the aggregated results on the client side.

This utility abstracts the complexities of interacting with homomorphic encryption and provides simple examples to help you integrate decryption into your workflow.

---

## Prerequisites

Before using **encproc-decryptor**, ensure that you have:
- A (your own or someone's else) running instance of the [encproc](https://github.com/collapsinghierarchy/encproc) engine.
- A valid JWT token for authentication with the encproc engine.
- [Go](https://golang.org/) installed if you wish to compile or modify the Go code.
- Basic knowledge of API requests.
---

## Installation

If you need to run or modify the Go code (e.g., `createStream.go`, `params.go`), make sure your Go environment is set up. Then install the dependencies:

```bash
go mod tidy
```

For client-side usage (HTML/JavaScript), simply open the HTML files in your browser or serve them via your preferred web server.

---

## Usage

### Creating a New Stream
By calling 
```bash
 go run createStream.go
``` 
everything will be set up correctly—provided you have a valid JWT token and the correct URL for the corresponding encproc engine. This program registers a new stream and generates a fresh public/secret key pair, which is stored in a file called `keypair.json` in the same directory. Store this file securely, as it also contains the stream ID. The format is:

```
{ 
        "id": "id",       
		"sk": "skBase64", 
		"pk": "pkBase64", 
}
```
Refer to the Go function in `createStream.go` to see how this is handled:
```GO
func createStream(apiURL, token string, publicKey []byte) (string, error)
```
---

### Contributing Data to a Stream

Once a stream is created, you can contribute encrypted data to it. A prerequisite is that you have the corresponding public key for the stream ID returned during the stream creation step. You can obtain the public key in two ways:

1. Load the public key from your `keypair.json` file.
2. Request the public key from the server by making a GET request to the following API endpoint:

```bash
GET /public-key/{id}
```

If the stream was created with an associated public key, the server will return the key that was provided during registration and stream creation. Once you have the public key, you are ready to encrypt the input.

To mark an integer for encryption, simply use `eng_push(integer)`, which fills up the input array in a stack-like fashion. When you are done filling up the stack of inputs, call `eng_enc(pk)` to trigger encryption of the entire stack.

#### Setting Up the WebAssembly Runtime

To use the `eng_push()` and `eng_enc(pk)` functions, you need to properly configure your WebAssembly (WASM) runtime environment for Go compilations. This involves several steps, similar to those implemented in `form.html`. We follow the instructions from the Go Wiki: [WebAssembly](https://go.dev/wiki/WebAssembly).

First, incorporate `wasm_exec.js` into your HTML. You can find this file in the repository, but you must ensure compatibility with our WASM binaries. The `wasm_exec.js` from this repository and from [encproc](https://github.com/collapsinghierarchy/encproc) are tested to work with the WASM binaries served by the encproc engine. Include it as follows:

```html
<script src="./wasm_exec.js"></script>
```

Next, fetch the appropriate Go-WASM binary, initialize the global `go` variable, and run the module:

```javascript
// Load and initialize the WASM module.
const go = new Go();
const wasmModule = await WebAssembly.instantiateStreaming(
  fetch("http://localhost:8000/static/encryption_module.wasm"),
  go.importObject
);
go.run(wasmModule.instance);
```

#### Encrypting and Contributing Data

Now, you can push values onto the input stack:

```javascript
// Call the WASM `eng_push` function for each input value.
error_msg_priv   = eng_push(priv);   // Push with the privacy preference.
error_msg_rating = eng_push(rating); // Push with the rating.
```

Encryption is handled by the `eng_enc(pk)` function, which returns a Base64-encoded string representing the ciphertext. See the following example from `form.html`:

```javascript
// Encrypt the data using the WASM `eng_encrypt` function.
const encryptedDataBase64 = eng_encrypt(publicKey);
```

Once encrypted, you are ready to contribute the data to the stream. In `form.html`, this is implemented as a simple POST request:

```javascript
// Prepare the payload.
const payload = JSON.stringify({ id, ct: encryptedDataBase64 });

// Send the encrypted data to the API.
const response = await fetch(`http://localhost:8000/contribute/aggregate/${id}`, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: payload,
});
```

This request forwards your encrypted data to the encproc engine, where it will be aggregated.

### Retrieving Aggregated Results

After enough contributions have been made, you can query the aggregated (encrypted) results. Use the same stream ID as before to fetch the results. In `results.md`, we provide an example of a simple GET request:

```javascript
const response = await fetch(`http://localhost:8000/snapshot/aggregate/${id}`, {
  method: "GET",
});
if (response.ok) {
  const data = await response.json();
  const { ct_aggr_byte_base64, id: responseID, sample_size } = data;
  // ...
}
```

The response includes:
- `id`: The stream ID from your request.
- `ct_aggr_byte_base64`: The aggregated ciphertext in Base64 format.
- `sample_size`: The number of contributions that were aggregated.

If the request fails, an error message will be returned.

To decrypt the aggregated ciphertext, you need the corresponding secret key; otherwise, the ciphertext is unusable. While decryption can be performed using the Lattigo Library, we provide a JavaScript decryption module—a Go decryption function compiled into WebAssembly (WASM). In our example, we load this WASM component from the encproc engine:

```javascript
const wasmFilePath = "http://localhost:8000/static/decrypt_results.wasm";
```

This module exports the function:

```javascript
// encryptedResults is the Base64-encoded ciphertext returned from the snapshot endpoint
// secretKey is the Base64-encoded secret key from keypair.json
const decryptedResults = decrypt_result(encryptedResults, secretKey);
```

The decrypted results will be an array of integers

## License

This repository is licensed under the [Apache-2.0 license](LICENSE).

---

## Contributing

Contributions to **encproc-decryptor** are welcome! Feel free to submit issues or pull requests with improvements or bug fixes. 

For major changes, please open an issue first to discuss what you would like to change.

---

## Additional Resources & Contact

- [encproc Repository](https://github.com/collapsinghierarchy/encproc) – The server-side engine that powers the encrypted processing as a service.
- Contact: encproc@gmail.com

Happy decrypting!
