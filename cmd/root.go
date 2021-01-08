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
	"log"
	"os"
	"os/exec"
	"strconv"

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
		if err != nil {
			fmt.Println(err.Error())
		}
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
