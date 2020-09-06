package tasks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
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

	cmdString := fmt.Sprintf("%s %s %s %s %s",
		"ffmpeg -y -i", payload.Src, payload.EncodeArgs, payload.Dst, "2>&1")

	cmd := exec.Command("sh", "-c",
		cmdString)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	// We're not using the -progress flag since it doesn't give us the duration
	// of the video which is important to determine the ETA. so we'll just parsing
	// the normal stdout.

	oneByte := make([]byte, 8)
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
		allRes += string(oneByte)
		getRatio(allRes)
	}

	err = cmd.Wait()

	log.Println("finished encoding")

	return err
}

var duration = 0
var allRes = ""
var lastPer = -1

func durToSec(dur string) (sec int) {
	durAry := strings.Split(dur, ":")
	if len(durAry) != 3 {
		return
	}
	hr, _ := strconv.Atoi(durAry[0])
	sec = hr * (60 * 60)
	min, _ := strconv.Atoi(durAry[1])
	sec += min * (60)
	second, _ := strconv.Atoi(durAry[2])
	sec += second
	return
}
func getRatio(res string) {
	durIdx := strings.Index(res, "Duration")
	// Checking if we've got a "Duration",
	// we need this so we can determine the ETA
	if durIdx >= 0 {

		dur := res[durIdx+10:]
		if len(dur) > 8 {
			dur = dur[0:8]

			duration = durToSec(dur)
			fmt.Printf("duration: %d (%s)", duration, res)
			allRes = ""
		}
	}
	// FFmpeg should give us a duration on startup,
	// so we kill here in the event that didn't happen.
	if duration == 0 {
		return
	}
	// From this point on it should be outputting normal encode stdout,
	// which we'll want to parse.
	timeIdx := strings.Index(res, "time=")
	// fpsIdx := strings.Index(res, "fps=")
	// sizeIdx := strings.Index(res, "size=")

	// strings.

	if timeIdx >= 0 {

		time := res[timeIdx+5:]
		if len(time) > 8 {
			time = time[0:8]
			sec := durToSec(time)
			per := (sec * 100) / duration
			if lastPer != per {
				lastPer = per
				fmt.Printf("Percentage: %d (%s)", per, res)
			}

			allRes = ""
		}
	}
}

func test() {
	cmdName := "ffmpeg -i 1.mp4  -acodec aac -vcodec libx264  cmd1.mp4 2>&1"
	cmd := exec.Command("sh", "-c", cmdName)
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	oneByte := make([]byte, 8)
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
		allRes += string(oneByte)
		getRatio(allRes)
	}
}
