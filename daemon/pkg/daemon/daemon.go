package daemon

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/svc"
	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/zerotier"
)

type NetworkDaemon struct {
	config Config
}

type Config struct {
	token     string
	networkID string
	nodeName  string
	domain    string
	namespace string
}

func NewConfig(token, networkID, nodeName, domain, namespace string) (Config, error) {
	return Config{token: token, networkID: networkID, nodeName: nodeName, domain: domain, namespace: namespace}, nil
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
	zt, err := zerotier.NewClient(d.config.token, d.config.networkID, d.config.nodeName, d.config.domain)
	if err != nil {
		return err
	}
	log.Info("daemon start")
	err = zt.Ensure()
	if err != nil {
		log.Error(err)
		zt.Stop()
		return err
	}
	time.Sleep(10 * time.Second)
	node, err := zt.GetNodeInfo()
	if err != nil {
		log.Error(err)
		zt.Stop()
		return err
	}
	if node.Domain != "" {
		err = svc.Ensure(node, d.config.namespace)
		if err != nil {
			log.Error(err)
		}
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
