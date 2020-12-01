package daemon

import (
	"time"

	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/zerotier"
)

type NetworkDaemon struct {
	config Config
}

type Config struct {
	token     string
	networkID string
}

func NewConfig(token, networkID string) (Config, error) {
	return Config{token: token, networkID: networkID}, nil
}

func NewDaemon(c Config) (NetworkDaemon, error) {
	return NetworkDaemon{}, nil
}

func (d NetworkDaemon) Start() error {
	return d.start()
}

func (d NetworkDaemon) start() error {
	zt, err := zerotier.NewClient(d.config.token, d.config.networkID)
	if err != nil {
		return err
	}
	for {
		err = zt.Ensure()
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second)
	}
}
