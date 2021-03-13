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
			Name:  "token, t",
			Usage: "Zerotier API Token",
		},
		cli.StringFlag{
			Name:  "networkID, n",
			Usage: "Zerotier Network ID",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "Zerotier Member name",
		},
		cli.StringFlag{
			Name:  "domain",
			Usage: "Zerotier Router domain",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "Zerotier Router namespace",
		},
	}
	app.Action = action
	app.Run(os.Args)
}

func action(c *cli.Context) error {
	token := c.String("token")
	if token == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	networkID := c.String("networkID")
	if networkID == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	nodeName := c.String("name")
	if nodeName == "" {
		cli.ShowAppHelp(c)
		return nil
	}
	domain := c.String("domain")
	namespace := c.String("namespace")

	config, err := daemon.NewConfig(token, networkID, nodeName, domain, namespace)
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
