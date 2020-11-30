package zerotier

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/takutakahashi/ztsvc-controller-daemon/pkg/node"
)

type Zerotier struct {
	executable Executable
}
type Executable interface {
	req(method, url string, params []byte) ([]byte, error)
	exec(cmd string) ([]byte, error)
}
type ZTExecutable struct {
	config Config
}

type Config struct {
	Token    string
	Endpoint string
}

func NewClient(token string) (Zerotier, error) {
	return Zerotier{
		executable: ZTExecutable{
			config: Config{
				Token:    token,
				Endpoint: "https://my.zerotier.com",
			},
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
func (zt Zerotier) updateMemberName(network, memberID, name string) error {
	p := struct {
		Name string `json:"name"`
	}{Name: name}
	params, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = zt.post("/api/network/"+network+"/member/"+memberID, params)
	return err
}
func (zt Zerotier) join(network string) error { return nil }
func (zt Zerotier) getNodeID() (string, error) {
	stdout, err := zt.executable.exec("info")
	if err != nil {
		return "", err
	}
	return strings.Split(string(stdout), " ")[2], nil
}
func (zt Zerotier) leave(network string) error { return nil }

func (e ZTExecutable) exec(cmd string) ([]byte, error) {
	return exec.Command("zerotier-cli", cmd).Output()
}
func (zt Zerotier) get(url string) ([]byte, error) {
	return zt.executable.req("GET", url, nil)
}
func (zt Zerotier) post(url string, params []byte) ([]byte, error) {
	return zt.executable.req("POST", url, params)
}
func (e ZTExecutable) req(method, url string, params []byte) ([]byte, error) {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	uri := e.config.Endpoint + url
	var req *http.Request
	var err error
	if method == "POST" {
		req, err = http.NewRequest(method, uri, bytes.NewBuffer(params))
	} else {
		req, err = http.NewRequest(method, uri, nil)
	}
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Authorization", "bearer "+e.config.Token)
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
