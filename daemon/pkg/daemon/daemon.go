package daemon

import (
	"time"

	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/zerotier"
)

type NetworkDaemon struct {
	config Config
}

type Config struct {
	token string
}

func NewConfig(token string) (Config, error) {
	return Config{token: token}, nil
}

func NewDaemon(c Config) (NetworkDaemon, error) {
	return NetworkDaemon{}, nil
}

func (d NetworkDaemon) Start() error {
	return d.start()
}

func (d NetworkDaemon) start() error {
	_, err := zerotier.NewClient(d.config.token)
	if err != nil {
		return err
	}
	for {
		time.Sleep(10 * time.Second)
	}
}
