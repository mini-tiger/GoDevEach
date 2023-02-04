package main

import (
	"context"
	"os"

	"github.com/apenella/go-ansible/pkg/playbook"
)

func main() {

	//ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
	//	Connection: "local",
	//	User:       "aleix",
	//}
	os.Chdir("/data/work/go/GoDevEach/ansible/ansibleplaybook-extravars-file")

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		//Inventory: "127.0.0.1,",
		Inventory: "hosts", // hostfile path
		ExtraVarsFile: []string{
			"@vars-file1.yml",
			"@vars-file2.yml",
		},
	}

	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks: []string{"site.yml"},
		//ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options: ansiblePlaybookOptions,
	}

	err := playbook.Run(context.TODO())
	if err != nil {
		panic(err)
	}
}
