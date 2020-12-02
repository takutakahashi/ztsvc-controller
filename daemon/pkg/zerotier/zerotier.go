package zerotier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type Zerotier struct {
	executable Executable
	config     Config
}
type Executable interface {
	req(method, url string, params []byte) ([]byte, error)
	exec(cmd []string) ([]byte, error)
}
type ZTExecutable struct {
	config Config
}

type Config struct {
	Token     string
	Endpoint  string
	NetworkID string
	NodeName  string
}

func NewClient(token, networkID, nodeName string) (Zerotier, error) {
	c := Config{
		Token:     token,
		Endpoint:  "https://my.zerotier.com",
		NetworkID: networkID,
		NodeName:  nodeName,
	}
	return Zerotier{
		config: c,
		executable: ZTExecutable{
			config: c,
		},
	}, nil
}

func (zt Zerotier) Ensure() error {
	err := zt.up()
	if err != nil {
		log.Error(err)
		return err
	}
	err = zt.join()
	if err != nil {
		log.Error(err)
		return err
	}
	node, err := zt.getNodeID()
	if err != nil {
		log.Error(err)
		return err
	}
	err = zt.authorize(node)
	if err != nil {
		log.Error(err)
		return err
	}
	err = zt.updateMemberName(node)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (zt Zerotier) Stop() error {
	return zt.leave()
}

func (zt Zerotier) getMembers(network string) (string, error) {

	b, err := zt.get("/api/network/" + network + "/member")
	return string(b), err
}
func (zt Zerotier) updateMemberName(memberID string) error {
	p := struct {
		Name string `json:"name"`
	}{Name: zt.config.NodeName}
	params, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = zt.post("/api/network/"+zt.config.NetworkID+"/member/"+memberID, params)
	return err
}

func (zt Zerotier) up() error {
	return exec.Command("/usr/sbin/zerotier-one", "-d").Start()
}
func (zt Zerotier) join() error {
	network := zt.config.NetworkID
	_, err := zt.executable.exec([]string{"join", network})
	return err
}

func (zt Zerotier) leave() error {
	network := zt.config.NetworkID
	_, err := zt.executable.exec([]string{"leave", network})
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
	_, err = zt.executable.req("POST", fmt.Sprintf("/api/network/%s/member/%s", zt.config.NetworkID, node), params)
	return err
}
func (zt Zerotier) getNodeID() (string, error) {
	stdout, err := zt.executable.exec([]string{"info"})
	if err != nil {
		return "", err
	}
	return strings.Split(string(stdout), " ")[2], nil
}

func (e ZTExecutable) exec(cmd []string) ([]byte, error) {
	log.Infof("exec command: %s", cmd)
	out, err := exec.Command("zerotier-cli", cmd...).Output()
	log.Info(string(out))
	log.Info(err)
	return out, err
}
func (zt Zerotier) get(url string) ([]byte, error) {
	return zt.executable.req("GET", url, nil)
}
func (zt Zerotier) post(url string, params []byte) ([]byte, error) {
	return zt.executable.req("POST", url, params)
}
func (e ZTExecutable) req(method, url string, params []byte) ([]byte, error) {
	log.Infof("req: %s", url)
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
	log.Infof("response status code: %d", resp.StatusCode)
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
