package daemon

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/zerotier"
)

type NetworkDaemon struct {
	config Config
}

type Config struct {
	token     string
	networkID string
	nodeName  string
}

func NewConfig(token, networkID, nodeName string) (Config, error) {
	return Config{token: token, networkID: networkID, nodeName: nodeName}, nil
}

func NewDaemon(c Config) (NetworkDaemon, error) {
	return NetworkDaemon{config: c}, nil
}

func (d NetworkDaemon) Start() error {
	return d.start()
}

func (d NetworkDaemon) start() error {
	zt, err := zerotier.NewClient(d.config.token, d.config.networkID, d.config.nodeName)
	if err != nil {
		return err
	}
	log.Info("daemon start")
	for {
		err = zt.Ensure()
		if err != nil {
			log.Error(err)
			return err
		}
		time.Sleep(10 * time.Second)
	}
}
