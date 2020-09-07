package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

	t := tasks.NewStore(NewCDN())

	taskserver.RegisterTasks(map[string]interface{}{
		"encode_video": t.EncodeVideo,
	})
	return taskserver
}

// NewCDN creates a connection to s3
func NewCDN() *s3.S3 {
	endpoint := os.Getenv("CDN_ENDPOINT")
	accessKeyID := os.Getenv("CDN_ACCESSKEYID")
	secretAccessKey := os.Getenv("CDN_SECRETACCESSKEY")

	// Configure to use CDN Server

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("ystv-wales-1"),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.New(s3Config)
	return s3.New(newSession)
}
