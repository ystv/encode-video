package worker

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/RichardKnop/machinery/v1"
)

// NewWorker creates a new worker
// not being used
func NewWorker() {
	cmd := exec.Command("ffmpeg", "-version")
	o, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			log.Fatalf("failed to find ffmpeg install")
		}
		log.Fatalf("failed to get ffmpeg version: %+v", err)
	}
	ver := strings.Split(string(o), " ")
	log.Println("encode-video: v0.1.0")
	log.Printf("using ffmpeg: v%s", ver[2])
}

// StartWorker creates a new worker
func StartWorker(taskserver *machinery.Server) error {
	worker := taskserver.NewWorker("machinery_worker", 10)
	if err := worker.Launch(); err != nil {
		return err
	}
	return nil
}
