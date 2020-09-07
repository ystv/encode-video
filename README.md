# encode-video

A go service, with both worker / server configurations available that handles batch video transcoding. Currently just a test project.

### Dependencies

- Redis (global)
- ffmpeg (worker)
- s3-like endpoint (worker)

### Server

`go run main.go server`

Runs a HTTP server listening on port `8082` accepting json encoded video encode requests at endpoint `/encode_video`.

### Worker

`go run main.go worker`

Will process video encode requests then copy the file over to the CDN.

### In action

Currently there is just the aforementioned HTTP endpoint to provide encode requests. For giving this project a whirl you can have an instance of a server and worker running, Redis and something S3-like, store a video file on it like what is in the example and run this curl command to operate:

`curl -d @single.json -X POST http://localhost:8082/encode_video`

There is also the option to do encode multiple files in one request with this query:

`curl -d @single.json -X POST http://localhost:8082/encode/multi`

### Endpoints

- `[POST] /encode/single` create a video encode request. Accepts payload object json and returns task UUID.
- `[POST] /encode/multi` create multiple encode requests. Accepts payload array json and returns first task UUID.
- `[GET] /pending` view jobs that haven't been accepted by a worker.
- `[GET] /status?uuid=[task uuid]` view task status.
