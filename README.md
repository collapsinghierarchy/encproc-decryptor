# encproc-decryptor


### Creating a New Stream
To create a new data stream:

```bash
curl -X POST http://localhost:8080/create-stream \
     -H "Authorization: Bearer <your_jwt_token>" \
     -d '{"pk": "<your_pk>"}'
```


### Contributing Data to a Stream

To contribute encrypted data to an existing stream:

```bash
curl -X POST http://localhost:8080/streams/{streamID}/contribute \
     -H "Authorization: Bearer <your_jwt_token>" \
     -d '{"encrypted_data": "<your_encrypted_data>"}'
```

### Retrieving Aggregated Results

To retrieve the aggregated result of a stream:

```bash
curl -X GET http://localhost:8080/streams/{streamID}/aggregate \
     -H "Authorization: Bearer <your_jwt_token>"
```