/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a host",
	Long:  `Add a new host to sshp.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			log.Fatal(err)
		}
		newhost.User = user

		promptHost := promptui.Prompt{
			Label:     "Host (or IP)",
			Templates: templates,
		}
		host, err := promptHost.Run()
		if err != nil {
			log.Fatal(err)
		}
		newhost.Host = host

		promptOwner := promptui.Prompt{
			Label:     "Owner",
			Templates: templates,
		}
		owner, err := promptOwner.Run()
		if err != nil {
			log.Fatal(err)
		}
		newhost.Owner = owner

		promptDesc := promptui.Prompt{
			Label:     "Description",
			Templates: templates,
		}
		desc, err := promptDesc.Run()
		if err != nil {
			log.Fatal(err)
		}
		newhost.Desc = desc

		prompt := promptui.Prompt{
			Label:     "Password",
			Templates: templates,
			Mask:      '*',
		}
		password, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		newhost.Password = password

		promptPort := promptui.Prompt{
			Label:     "Port (22)",
			Templates: templates,
		}
		port, err := promptPort.Run()
		if err != nil {
			log.Fatal(err)
		}
		if port == "" {
			port = "22"
		}
		newhost.Port, err = strconv.Atoi(port)
		if err != nil {
			log.Fatal(err)
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
	},
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

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
