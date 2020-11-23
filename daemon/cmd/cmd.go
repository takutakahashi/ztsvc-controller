package main

import (
	"log"
	"os"

	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/daemon"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "ztsvc controller daemon"
	app.Usage = "daemon process for ztsvc"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config.yaml filepath",
		},
	}
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) error {
	configPath := c.String("config")
	if configPath == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	config, err := daemon.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	d, err := daemon.NewDaemon(config)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	d.Start()
	return nil
}
