# encode-video

A go service, with both worker / server configurations available that handles batch video transcoding. Currently just a test project.

## Dependencies

- Redis (global)
- ffmpeg (worker)
- postgres (at some point)
- s3-like endpoint (at some point)

## Server `go run main.go server`

Runs a HTTP server listening on port `8082` accepting json encoded video encode requests at endpoint `/encode_video`.

## Worker `go run main.go worker`

Will process video encode requests then copy the file over to the CDN.
