/*
Copyright © 2020 Spencer Dupre <spencer.dupre@gmail.com>

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
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	portFlag     int
	dontCopyFlag bool
	// torFlag bool
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a host",
	Long:  `Add a new host to sshp.`,
	Run: func(cmd *cobra.Command, args []string) {
		newhost := Host{}
		newhost.Port = portFlag
		newhost.Tor = false
		newhost.Timestamp = time.Now().Unix()

		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Success: "{{ . | bold }} ",
		}

		promptUser := promptui.Prompt{
			Label:     "Username",
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

		if !dontCopyFlag {
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
			authhost(newhost)
		}

		hosts, err := gethosts(HostsFile)
		if err != nil {
			log.Fatal(err)
		}
		newhost.Password = ""
		hosts = append(hosts, newhost)
		err = writehosts(hosts)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func authhost(host Host) {
	hostname := host.Host
	port := host.Port
	user := host.User
	pass := host.Password

	var pubkeyfile = filepath.Join(home, ".ssh", "id_rsa.pub")
	pubkey, err := ioutil.ReadFile(pubkeyfile)
	if err != nil {
		log.Fatal(err)
	}

	// The bash command we will run on the remote machine, to copy our key
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
	var client *ssh.Client
	if torFlag := false; torFlag {
		client, err = proxiedSSHClient("localhost:9050", fullhost, config)
	} else {
		client, err = ssh.Dial("tcp", fullhost, config)
	}
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

	// connect to stdout and stderr
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// run single command
	err = sess.Run(cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().IntVarP(&portFlag, "port", "p", 22, "Port to use for SSH connection")
	addCmd.Flags().BoolVarP(&dontCopyFlag, "dont-copy", "x", false, "Don't copy public key to remote host")
	// addCmd.Flags().BoolVarP(&torFlag, "tor", "t", false, "Use tor when connecting to this host")
}
