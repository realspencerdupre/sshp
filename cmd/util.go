package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
)

type Host struct {
	Host      string
	User      string
	Desc      string
	Timestamp int
	Owner     string
	Password  string
	Port      int
}

var home, _ = homedir.Dir()
var HostsFile = filepath.Join(home, ".sshp_hosts.json")

func gethosts(path string) ([]Host, error) {
	var hosts []Host
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

func selecthost(hosts []Host) (Host, error) {
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
		return Host{}, err
	}

	return hosts[i], nil
}
