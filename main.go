package main

import (
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/urfave/cli"
	"github.com/ystv/encode-video/server"
	"github.com/ystv/encode-video/utils"
	"github.com/ystv/encode-video/worker"
)

var (
	app        *cli.App
	taskserver *machinery.Server
)

func init() {
	app = cli.NewApp()
	taskserver = utils.NewMachineryServer()
}

func main() {
	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the encode server",
			Action: func(c *cli.Context) {
				server.StartServer(taskserver)
			},
		},
		{
			Name:  "worker",
			Usage: "Run the encode worker",
			Action: func(c *cli.Context) {
				worker.StartWorker(taskserver)
			},
		},
	}
	app.Run(os.Args)
}
