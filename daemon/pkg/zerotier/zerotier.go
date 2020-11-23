package zerotier

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/node"
)

type Zerotier struct {
	config Config
}
type Config struct {
	Token    string
	Endpoint string
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
func (zt Zerotier) get(url string) ([]byte, error) {
	return zt.req("GET", url, struct{}{})
}
func (zt Zerotier) post(url string, params struct{}) ([]byte, error) {
	return zt.req("POST", url, params)
}
func (zt Zerotier) req(method, url string, params struct{}) ([]byte, error) {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(method, zt.config.Endpoint+url, bytes.NewBuffer(paramsJSON))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Authorization", "Bearer "+zt.config.Token)
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
