package utils

import (
	"fmt"
	"log"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/ystv/encode-video/tasks"
)

// NewMachineryServer establishes our connection
func NewMachineryServer() *machinery.Server {
	taskserver, err := machinery.NewServer(&config.Config{
		Broker:        "redis://localhost:6379",
		ResultBackend: "redis://localhost:6379",
	})
	if err != nil {
		log.Fatalf("%+v", fmt.Errorf("failed to connect to broker: %w", err))
	}
	taskserver.RegisterTasks(map[string]interface{}{
		"encode_video": tasks.EncodeVideo,
	})
	return taskserver
}
