package tasks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type (
	// Payload the object that is used to communicate in the MQ
	Payload struct {
		Src        string `json:"src"`        // Location of source file on CDN
		Dst        string `json:"dst"`        // Destination of finished encode on CDN
		EncodeName string `json:"encodeName"` // Here for pretty logging
		EncodeArgs string `json:"encodeArgs"` // Encode arguments
	}
	// Stats represents statistics on the current encode job
	Stats struct {
		Duration   int    `json:"duration"`
		Percentage int    `json:"percentage"`
		Frame      int    `json:"frame"`
		FPS        int    `json:"fps"`
		Bitrate    string `json:"bitrate"`
		Size       string `json:"size"`
		Time       string `json:"time"`
	}
)

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

	bytes := make([]byte, 100)
	stats := &Stats{}
	allRes := ""
	start := time.Now()
	for {
		_, err := stdout.Read(bytes)
		if err != nil {
			err = fmt.Errorf("failed to read stdout: %w", err)
			break
		}
		allRes += string(bytes)
		ok := getStats(allRes, stats)
		if ok {
			allRes = ""
			log.Printf("%+v", stats)
		}
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	log.Printf("finished encoding %s/%s - completed in %s", payload.Src, payload.EncodeName, time.Since(start))

	return nil
}

func durToSec(dur string) (sec int) {
	// So we're kind of limiting our videos to 24h which isn't ideal
	// shouldn't crash the application hopefully XD
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
func getStats(res string, s *Stats) bool {

	durIdx := strings.Index(res, "Duration")
	// Checking if we've got a "Duration",
	// we need this so we can determine the ETA
	if durIdx >= 0 {

		dur := res[durIdx+10:]
		if len(dur) > 8 {
			dur = dur[0:8]

			s.Duration = durToSec(dur)
			return true
		}
	}
	// FFmpeg should give us a duration on startup,
	// so we kill here in the event that didn't happen.
	if s.Duration == 0 {
		return false
	}

	frameIdx := strings.LastIndex(res, "frame=")
	fpsIdx := strings.LastIndex(res, "fps=")
	bitrateIdx := strings.LastIndex(res, "bitrate=")
	sizeIdx := strings.LastIndex(res, "size=")
	timeIdx := strings.Index(res, "time=")

	if timeIdx >= 0 {
		// From this point on it should be outputting normal encode stdout,
		// which we'll want to parse.

		frame := strings.Fields(res[frameIdx+6:])
		fps := strings.Fields(res[fpsIdx+4:])
		bitrate := strings.Fields(res[bitrateIdx+8:])
		size := strings.Fields(res[sizeIdx+5:])
		time := res[timeIdx+5:]

		if len(time) > 8 {
			time = time[0:8]
			sec := durToSec(time)
			per := (sec * 100) / s.Duration
			if s.Percentage != per {
				s.Percentage = per
				// Just doing to reuse this int variable for each item
				integer, _ := strconv.Atoi(frame[0])
				s.Frame = integer
				integer, _ = strconv.Atoi(fps[0])
				s.FPS = integer
				s.Bitrate = bitrate[0]
				s.Size = size[0]
				s.Time = time
			}
			return true
		}
	}
	return false
}
