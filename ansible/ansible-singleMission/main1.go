package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/apenella/go-ansible/ansible"
)

func main() {
	runner := ansible.NewRunner()

	config := &ansible.AnsibleRunnerConfig{
		Inventory:         "/path/to/inventory",
		Limit:             "all",
		ModulePath:        "/path/to/module",
		User:              "your-username",
		PrivateKeyFile:    "/path/to/private/key/file",
		ConnectionOptions: []string{"ssh_args=-o ForwardAgent=yes"},
	}
	runner.Configure(config)

	command := &ansible.AnsibleAdhocCommand{
		Patterns:   []string{"all"},
		ModuleName: "shell",
		ModuleArgs: "ls -l /etc/",
	}
	result, err := runner.RunAdhocCommand(command)
	if err != nil {
		log.Fatal(err)
	}

	var response []ansible.AdhocResponse
	if err := json.Unmarshal(result.StdoutBytes, &response); err != nil {
		log.Fatal(err)
	}

	for _, r := range response {
		fmt.Println(r.Host)
		fmt.Println(r.Task.Name)
		fmt.Println(r.Task.Description)
		fmt.Println(r.Task.Status)
		fmt.Println(r.Task.StartTime)
		fmt.Println(r.Task.EndTime)
		fmt.Println(r.Task.Duration)
		fmt.Println(r.Task.Stdout)
		fmt.Println(r.Task.Stderr)
	}
}
