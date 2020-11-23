package zerotier

import "github.com/takutakahashi/ztsvc-controller-daemon/pkg/node"

type Zerotier struct{}

func NewClient() (Zerotier, error) {
	return Zerotier{}, nil
}

func (zt Zerotier) Ensure(node.Node) error {

	return nil
}
