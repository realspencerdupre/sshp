package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/manifoldco/promptui"
)

type Host struct {
	UserHost  string
	Desc      string
	Timestamp int
	Owner     string
	Password  string
	Port      int
}

func main() {
	data, err := ioutil.ReadFile("./hosts.json")
	if err != nil {
		fmt.Print(err)
	}

	var hosts []Host

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		fmt.Println("error:", err)
	}

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
		return
	}

	fmt.Printf("Connecting to: %s\n", hosts[i].Desc)
}
