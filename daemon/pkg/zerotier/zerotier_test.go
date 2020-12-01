package zerotier

import (
	"errors"
	"testing"

	"github.com/labstack/gommon/log"
)

type ExecutableMock struct{}

func (e ExecutableMock) req(method, url string, params []byte) ([]byte, error) {
	return []byte{}, nil
}
func (e ExecutableMock) exec(cmd []string) ([]byte, error) {
	log.Info(cmd)
	switch cmd[0] {
	case "info":
		return []byte("200 info 0000000000 1.4.6 ONLINE"), nil
	case "join":
		return []byte("OK"), nil
	default:
		return nil, errors.New("error")
	}
}
func Mock() Zerotier {
	network := "0000000000000000"
	nodeName := "node"

	zt, _ := NewClient("", network, nodeName)
	return zt
}

func TestEnsure(t *testing.T) {
	zt := Mock()
	zt.executable = ExecutableMock{}
	err := zt.Ensure()
	if err != nil {
		t.Fatal(err)
	}
}
