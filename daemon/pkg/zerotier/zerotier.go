package zerotier

import "github.com/takutakahashi/ztsvc-controller-daemon/pkg/node"

type Zerotier struct {
	config Config
}
type Config struct {
	Token string
}

func NewClient(token string) (Zerotier, error) {
	return Zerotier{
		config: Config{
			Token: token,
		},
	}, nil
}

func (zt Zerotier) Ensure(node.Node) error {

	return nil
}

func (zt Zerotier) getNetwork()   {}
func (zt Zerotier) getMembers()   {}
func (zt Zerotier) addMember()    {}
func (zt Zerotier) deleteMember() {}
func (zt Zerotier) addVip()       {}
func (zt Zerotier) deleteVip()    {}
