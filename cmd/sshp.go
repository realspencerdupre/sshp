package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
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

func addhost() error {
	newhost := Host{}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Success: "{{ . | bold }} ",
	}

	promptUser := promptui.Prompt{
		Label:     "User",
		Templates: templates,
	}
	user, err := promptUser.Run()
	if err != nil {
		return err
	}
	newhost.User = user

	promptHost := promptui.Prompt{
		Label:     "Host (or IP)",
		Templates: templates,
	}
	host, err := promptHost.Run()
	if err != nil {
		return err
	}
	newhost.Host = host

	promptOwner := promptui.Prompt{
		Label:     "Owner",
		Templates: templates,
	}
	owner, err := promptOwner.Run()
	if err != nil {
		return err
	}
	newhost.Owner = owner

	promptDesc := promptui.Prompt{
		Label:     "Description",
		Templates: templates,
	}
	desc, err := promptDesc.Run()
	if err != nil {
		return err
	}
	newhost.Desc = desc

	prompt := promptui.Prompt{
		Label:     "Password",
		Templates: templates,
		Mask:      '*',
	}
	password, err := prompt.Run()
	if err != nil {
		return err
	}
	newhost.Password = password

	promptPort := promptui.Prompt{
		Label:     "Port (22)",
		Templates: templates,
	}
	port, err := promptPort.Run()
	if err != nil {
		return err
	}
	if port == "" {
		port = "22"
	}
	newhost.Port, err = strconv.Atoi(port)
	if err != nil {
		return err
	}

	hosts, err := gethosts(HostsFile)
	if err != nil {
		log.Fatal(err)
	}
	hosts = append(hosts, newhost)
	err = writehosts(hosts)
	if err != nil {
		log.Fatal(err)
	}
	authhost(newhost)
	return nil
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

func rmhost() error {
	hosts, err := gethosts(HostsFile)
	if err != nil {
		log.Fatal(err)
	}
	selectedhost, err := selecthost(hosts)
	if err != nil {
		log.Fatal(err)
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Success: "{{ . | bold }} ",
	}
	promptconfirm := promptui.Prompt{
		Label:     "Are you sure? (y)",
		Templates: templates,
	}
	confirm, err := promptconfirm.Run()
	if err != nil {
		return err
	}
	confirm = strings.ToLower(confirm)
	if !(strings.HasPrefix(confirm, "y")) {
		fmt.Println("Do nothing")
		return nil
	}

	var newhosts []Host
	for i := range hosts {
		if hosts[i] != selectedhost {
			newhosts = append(newhosts, hosts[i])
		}
	}
	fmt.Println("Removing", selectedhost.Desc)
	err = writehosts(newhosts)
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

func connecthost(host Host) error {
	fullhost := fmt.Sprintf("%s@%s", host.User, host.Host)
	cmd := exec.Command("ssh", "-p", strconv.Itoa(host.Port), fullhost)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func authhost(host Host) error {
	hostname := host.Host
	port := host.Port
	user := host.User
	pass := host.Password

	var pubkeyfile = filepath.Join(home, ".ssh", "id_rsa.pub")
	pubkey, err := ioutil.ReadFile(pubkeyfile)
	if err != nil {
		log.Fatal(err)
	}

	cmd := fmt.Sprintf("echo '%s' >> $HOME/.ssh/authorized_keys", pubkey)

	// ssh client config
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect
	fullhost := fmt.Sprintf("%s:%d", hostname, port)
	client, err := ssh.Dial("tcp", fullhost, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// start session
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	// setup standard out and error
	// uses writer interface
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// run single command
	err = sess.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
