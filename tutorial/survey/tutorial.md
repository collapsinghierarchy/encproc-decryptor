# Introduction & Motivation
This tutorial will describe how to setup your own encrypted survey webapp leveraging the API engine.

## Prior Knowledge
The main motivation for the API engine was to reduce the required cryptographical expertise threshold for developers to being able to implement or enhance web applications with homomorphic encryption. That being said the least requirement for the target audience of this tutorial is to understand what web development means. You should be at least familiar with terminologies such as API, Endpoint, GET/POST or some other HTTP(S) Requests, URL, Javascript, HTML and so on. Also some basic understanding of security concepts would be very helpful in understanding why and what you are doing here. 

>The good part is you dont need to understand how homomorphic encryption actually works. Only that it's an asymmetric encryption scheme and that there are public and secret keys involved, where the public key is public and the secret key should remain secret **and ideally the secret key is not in the same hands as the ciphertexts you are sending to. (Otherwise you security guarantees are void and the webapp is useless)** 

Have fun!

## Structure, Process and Risks of a General Survey WebApp
The following diagram is an example of a regular generic survey webapp (without any encryption) that is deployed within the same organisation as the survey participants and survey creators belong to. 
![alt text](image-2.png)
The legend is the following:
- Circle + Text --> Processing steps.
- <font color="#add8e6">BLUE: Personal Data / Personally Identifiable Data (PII) of various sensitivity levels that can be linked to the survey participant, i.e. the **data subject**</font>
- <font color="green">GREEN: Parties that are authorized to interact with the personal data of our data subjects.</font>
- <font color="#a13684">VIOLET: "Anonymous" data, i.e. not linkable to the **data subject** (irreversibly)</font> (<font color="gray">Currently under heavy debate on how to prove that the presumably anonymous data is (A) not linkable and (B) irreversible)</font>
- <font color="red">RED: Parties that are unauthorized to interact with the personal data of our data subjects.</font>
- <font color="#af593e">BROWN: Encrypted (or according to GDPR -- Pseudonymous) Data, which is most of the time TLS-traffic. But later in the tutorial it will spread.</font>

The overall process is simple and straight forward.
1. A survey participant requests the survey form and contributes (HTTPS) answers.
2. A survey creators decides that time is over and it is time to get the survey results, i.e., compute the aggregation.
3. The Back-End computes the aggregation. 
4. The survey creators gets the results of the aggregation.

The obvious risks/attack vectors are:
- Malicious Intrusion: Malware that leaks the data or maybe even ask for ransom prior to leaking the data.
- Internal Adversary: Regular employee abusing access.
- Human Error: Admin or other authorized employees unintentionally leaking access for a public actor due to human error.

Not so obvious may be:
- Output Privacy: Eventually (if the results returned to survey creator are not all aggregated or the sample size is too small), the survey creator may constitute a fourth risk/attack vector. 

In this tutorial we will see how we can address the first 3 risks, i.e., Malicious Intrusion, Interal Adversary, Human Error.

##  Structure, Process and Risks of a General Encrypted Survey WebApp
As you can see the changes of the overall processing are minimal. There is only the additional step of the key generation within the Survey Creation Front-End component.
![alt text](image-1.png)
The main component is now that the Survey Front-End Component at the **Data Subject** side encrypts the answers **homomorphically** before contributing them to the Back-End. This allows the Back-End to aggregate the encrypted answers, i.e. **ciphertexts/ciphtxts** or as GDPR practioneers call them **pseudonyms**, in the same way as it was being done before (albeit now we do not compute the "simple" arithmetic **+**, but the **homomorphic encryption +**, which is from a cryptologists perspective one of the most trivial things among the standard **+**).

So how does this tackle the 3 previously mentioned attack vectors?
- Malicious Intrusion: The malware learns only the ciphertexts, while the secret key remains out of reach on the end-device of the survey creator. (Ideally it is deleted after the decryption) 
- Internal Adversary: Same as the first category.
- Human Error: Same as the first category. 

This is due to the **honest-but-curious** attacker model, where it is assumed that any actor gaining access to the encrypted data does so only passively and does not try to manipulate the data. Protections against an **active** adversary are *somewhat* possible, but out-of-scope for this tutorial.

# So How to Create the Encrypted Survey WebApp?
First of all we need to understand how the API engine [encproc](https://github.com/collapsinghierarchy/encproc) works. The construction of the encrypted survey webapp will follow then straight forwardly because it's simply one of the most simple use-cases for the engine and the homomorphic encryption in general. You can find a ready-to-use example in the [demo](https://github.com/collapsinghierarchy/encproc-decryptor/tree/main/survey%20demo) repository.

If you don't want to use the engine, then you can try to use the homomorphic encryption [lattigo](https://github.com/tuneinsight/lattigo) directly and build your own server. The API engine does not bring nothing special to the table. It is a simple wrapper around the lattigo library.

## The API Engine EncProc
In our survey webApp the engine will take the role of the **Survey Back-End** that is self-hosted within the organisation. 

> **Note:** Alternative architectures may be plenty. For example you can use your already existing **not-encrypted** survey webApp and enhance it with the API engine for aggregating the ciphertexts. This will have the positive side-effect of **data minimization** and **separation of duties** such that the encrypted data processing API engine does not receive any client-linkable meta-data of the survey participants.

### Setup EncProc
For the sake of simplicity of this tutorial we refer to the general setup of the engine to the [README](https://github.com/collapsinghierarchy/encproc) of the according repository. For this tutorial we will be using an already running instance of encproc ``http://217.154.80.44:8080``. The survey [demo](https://github.com/collapsinghierarchy/encproc-decryptor/tree/main/survey%20demo) example is already pre-configured to use this instance. 
You can try to access the URL from your browser and you will receive ``heyho`` back in case it's up and running.

### Communicate with EncProc

The essential API endpoints for the encrypted webApp are the following:

- **`POST /create-stream`**:  
  Create a new data stream. This endpoint requires a valid JWT token and a properly generated public key. It returns the registered stream ID.

- **`POST /contribute/aggregate/{id}`**:  
  Contribute encrypted data to a specific stream identified by `{id}`.

- **`GET /snapshot/aggregate/{id}`**:  
  Retrieve the aggregated result of a specific stream identified by `{id}`.

- **`GET /public-key/{id}`**:  
  Retrieve the public key associated with a specific stream identified by `{id}` for encryption purposes.

## The Encrypted WebApp
Bringing back the picture from above:
![alt text](image-1.png)
You can see the sequence of functions that we will need to implement on the client-side of a <font color="#add8e6">Survey Participant</font> and a <font color="green">Survey Creator</font>. The functional requirements are therefore structured as follows:
- The survey creator will need to
    - **SC_REQ1** Generate a keypair.
    - **SC_REQ2** Create a survey form.
    - **SC_REQ3** Send the public key to the back-end.
    - **SC_REQ4** Request & receive the encrypted aggregated results.
    - **SC_REQ5** Decrypt the aggregated results.
    
- The survey participant will need to
    - **SP_REQ1** Request and fill our the survey form.
    - **SP_REQ2** Receive the according public key.
    - **SP_REQ3** Encrypt and send the answers.

- The survey back-end will need to
    - to do a lot of stuff about which we do not care about in this tutorial. If you are interested you can look into the [encproc implementation](https://github.com/collapsinghierarchy/encproc)

In order to implemented a **minimal** webapp that meets all of these requirements we will need to implement several components.

> This is only a minimal version of such a webapp that fits a simple tutorial. I hope that the reader understands that there are many more aspects that needs to be implemented for a user friendly survey web app. My hope is that people do exactly this after the tutorial ^^

- The survey creator front-end
    - HTML+Javascript
- The survey participant front-End
    - HTML+Javascript

### Survey Creator Front-End

### Survey Participant Front-End