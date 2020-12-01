package daemon

import (
	"os"
	"os/signal"
	"syscall"

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
	sigs := make(chan os.Signal, 1)
	done := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	zt, err := zerotier.NewClient(d.config.token, d.config.networkID, d.config.nodeName)
	if err != nil {
		return err
	}
	log.Info("daemon start")
	err = zt.Ensure()
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		sig := <-sigs
		log.Info(sig)
		log.Info("terminating")
		err := zt.Stop()
		if err != nil {
			log.Errorf("can't leave: %s", err)
		}
		done <- err
	}()
	return <-done
}
