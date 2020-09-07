# encode-video

A go service, with both worker / server configurations available that handles batch video transcoding. Currently just a test project.

### Dependencies

- Redis (global)
- ffmpeg (worker)
- postgres (at some point)
- s3-like endpoint (at some point)

### Server `go run main.go server`

Runs a HTTP server listening on port `8082` accepting json encoded video encode requests at endpoint `/encode_video`.

### Worker `go run main.go worker`

Will process video encode requests then copy the file over to the CDN.

### In action

Currently there is just the forementioned HTTP endpoint to provide encode requests. For giving this project a whirl you can have an instance of a server and worker running then pop a file called `source.mp4` and run this curl command to operate:

`curl -d @example.json -X POST http://localhost:8082/encode_video`

You can change the request there to match the file. Later versions will hopefully pull from the CDN first.
