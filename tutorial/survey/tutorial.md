# Introduction & Motivation
This tutorial will describe how to setup your own encrypted survey webapp leveraging the API engine.

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

# So how to create the encrypted survey WebApp?
