package zerotier

import (
	"errors"
	"testing"
)

type ExecutableMock struct{}

func (e ExecutableMock) req(method, url string, params []byte) ([]byte, error) {
	return []byte{}, nil
}
func (e ExecutableMock) exec(cmd string) ([]byte, error) {
	switch cmd {
	case "info":
		return []byte("200 info 0000000000 1.4.6 ONLINE"), nil
	default:
		return nil, errors.New("error")
	}
}

func TestRequest(t *testing.T) {
	network := "0000000000000000"
	node := "0000000000"
	zt, err := NewClient("")
	zt.executable = ExecutableMock{}
	if err != nil {
		t.Fatal("error", err)
	}
	err = zt.join(network)
	if err != nil {
		t.Fatal(err)
	}
	nodeID, err := zt.getNodeID()
	if err != nil || nodeID != node {
		t.Fatalf("expected: %s, actual: %s", node, nodeID)
	}
	err = zt.updateMemberName(network, node, "node01")
	if err != nil {
		t.Fatal("error", err)
	}
}
