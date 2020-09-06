package tasks

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

// Payload the object that is used to communicate in the MQ
type Payload struct {
	Src        string `json:"src"`        // Location of source file on CDN
	Dst        string `json:"dst"`        // Destination of finished encode on CDN
	EncodeName string `json:"encodeName"` // Here for pretty logging
	EncodeArgs string `json:"encodeArgs"` // Encode arguments
}

// DecodeToTask converts the b64 encoded task and converts it to the payload object
func DecodeToTask(msg string, task interface{}) (err error) {
	decodedstg, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return
	}
	msgBytes := []byte(decodedstg)
	err = json.Unmarshal(msgBytes, task)
	if err != nil {
		return
	}
	return
}

// EncodeVideo encodes a video
func EncodeVideo(b64payload string) error {
	payload := Payload{}
	err := DecodeToTask(b64payload, &payload)

	log.Printf("encode video: %s/%s", payload.Src, payload.EncodeName)
	return err
}
