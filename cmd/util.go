package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

type Host struct {
	Host      string
	User      string
	Desc      string
	Timestamp int
	Owner     string
	Password  string
	Port      int
	Tor       bool
}

var home, _ = homedir.Dir()
var HostsFile = filepath.Join(home, ".sshp_hosts.json")

func gethosts(path string) ([]Host, error) {
	var hosts []Host

	// Touch hosts file with empty array if it doesn't exist
	_, err := os.Stat(HostsFile)
	if os.IsNotExist(err) {
		file, err := os.Create(HostsFile)
		if err != nil {
			log.Fatal(err)
		}
		output, err := json.Marshal(hosts)
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.Write(output)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	// Read the hosts file
	data, err := ioutil.ReadFile(HostsFile)
	if err != nil && strings.HasSuffix(err.Error(), "such file or directory") {
		return hosts, errors.New("No hosts configured")
	}
	err = json.Unmarshal(data, &hosts)
	if err != nil {
		fmt.Println("error:", err)
	}
	return hosts, nil
}

func writehosts(hosts []Host) error {
	output, err := json.Marshal(hosts)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(HostsFile, output, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func selecthost(hosts []Host) (int, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .UserHost }}?",
		Active:   "> {{ .Desc | cyan }} ({{ .Owner }})",
		Inactive: "  {{ .Desc | cyan }} ({{ .Owner }})",
		Selected: "> {{ .Desc | red | cyan }}",
	}
	prompt := promptui.Select{
		Label:     "Select Day",
		Items:     hosts,
		Templates: templates,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1, err
	}

	return i, nil
}

func proxiedSSHClient(proxyAddress string, sshServerAddress string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}

	conn, err := dialer.Dial("tcp", sshServerAddress)
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddress, sshConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}
