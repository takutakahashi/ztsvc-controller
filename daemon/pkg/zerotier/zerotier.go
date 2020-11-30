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
			Token:    token,
			Endpoint: "https://my.zerotier.com",
		},
	}, nil
}

func (zt Zerotier) Ensure(node.Node) error {

	return nil
}

func (zt Zerotier) getMembers(network string) (string, error) {

	b, err := zt.get("/api/network/" + network + "/member")
	return string(b), err
}
func (zt Zerotier) register()  {}
func (zt Zerotier) leave()     {}
func (zt Zerotier) addVip()    {}
func (zt Zerotier) deleteVip() {}
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
	uri := zt.config.Endpoint + url
	var req *http.Request
	if method == "POST" {
		p := bytes.NewBuffer(paramsJSON)
		req, err = http.NewRequest(method, uri, p)
	} else {
		req, err = http.NewRequest(method, uri, nil)
	}
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Authorization", "bearer "+zt.config.Token)
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
