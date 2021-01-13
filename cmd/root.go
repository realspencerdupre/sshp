/*
Copyright © 2021 Spencer Dupre <spencer.dupre@gmail.com>

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
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "sshp",
	Short: "SSH login manager",
	Long:  `sshp is a quick way to store many ssh logins in one, easy to reach place.`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := gethosts(HostsFile)
		if err != nil {
			log.Fatal(err)
		}
		hostIndex, err := selecthost(hosts)

		// host timestamp is set to current time,
		// so we have the most recently connected host
		hosts[hostIndex].Timestamp = time.Now().Unix()
		writehosts(hosts)

		if err != nil {
			fmt.Println(err.Error())
		}

		// Full hostname, like user@example.com
		fullhost := fmt.Sprintf(
			"%s@%s",
			hosts[hostIndex].User,
			hosts[hostIndex].Host,
		)
		command := []string{
			"ssh",
			"-p",
			strconv.Itoa(hosts[hostIndex].Port),
			fullhost,
		}

		// Start the session until the user ends it
		session := exec.Command(command[0], command[1:]...)
		session.Stdout = os.Stdout
		session.Stdin = os.Stdin
		session.Stderr = os.Stderr
		err = session.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
