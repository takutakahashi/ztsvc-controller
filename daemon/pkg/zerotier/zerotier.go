package zerotier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Zerotier struct {
	executable Executable
	config     Config
}
type Executable interface {
	req(method, url string, params []byte) ([]byte, error)
	exec(cmd string) ([]byte, error)
}
type ZTExecutable struct {
	config Config
}

type Config struct {
	Token     string
	Endpoint  string
	NetworkID string
}

func NewClient(token, networkID string) (Zerotier, error) {
	c := Config{
		Token:     token,
		Endpoint:  "https://my.zerotier.com",
		NetworkID: networkID,
	}
	return Zerotier{
		config: c,
		executable: ZTExecutable{
			config: c,
		},
	}, nil
}

func (zt Zerotier) Ensure() error {
	err := zt.join()
	if err != nil {
		return err
	}
	node, err := zt.getNodeID()
	if err != nil {
		return err
	}
	err = zt.authorize(node)
	if err != nil {
		return err
	}
	err = zt.updateMemberName(node, os.Getenv("HOSTNAME"))
	if err != nil {
		return err
	}
	return nil
}

func (zt Zerotier) getMembers(network string) (string, error) {

	b, err := zt.get("/api/network/" + network + "/member")
	return string(b), err
}
func (zt Zerotier) updateMemberName(memberID, name string) error {
	p := struct {
		Name string `json:"name"`
	}{Name: name}
	params, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = zt.post("/api/network/"+zt.config.NetworkID+"/member/"+memberID, params)
	return err
}
func (zt Zerotier) join() error {
	network := zt.config.NetworkID
	_, err := zt.executable.exec(fmt.Sprintf("join %s", network))
	return err
}

func (zt Zerotier) leave() error {
	network := zt.config.NetworkID
	_, err := zt.executable.exec(fmt.Sprintf("leave %s", network))
	return err
}

func (zt Zerotier) authorize(node string) error {
	type c struct {
		Authorized bool `json:"authorized"`
	}
	p := struct {
		Config c
	}{Config: c{Authorized: true}}
	params, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = zt.executable.req("POST", fmt.Sprintf("/network/%s/member/%s", zt.config.NetworkID, node), params)
	return err
}
func (zt Zerotier) getNodeID() (string, error) {
	stdout, err := zt.executable.exec("info")
	if err != nil {
		return "", err
	}
	return strings.Split(string(stdout), " ")[2], nil
}

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
