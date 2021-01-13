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
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a host",
	Long:  `Select and remove a host from the list.`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts, err := gethosts(HostsFile)
		if err != nil {
			log.Fatal(err)
		}
		selectedHostIndex, err := selecthost(hosts)
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
			log.Fatal(err)
		}
		confirm = strings.ToLower(confirm)
		if !(strings.HasPrefix(confirm, "y")) {
			fmt.Println("Do nothing")
			log.Fatal(err)
		}

		var newhosts []Host
		for i := range hosts {
			if i != selectedHostIndex {
				newhosts = append(newhosts, hosts[i])
			}
		}
		fmt.Println("Removing", hosts[selectedHostIndex].Desc)
		err = writehosts(newhosts)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
