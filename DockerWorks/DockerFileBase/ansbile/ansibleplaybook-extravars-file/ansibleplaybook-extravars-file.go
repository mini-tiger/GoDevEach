package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
	"io"
	"os"

	"github.com/apenella/go-ansible/pkg/playbook"
)

func main() {

	//ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
	//	Connection: "local",
	//	User:       "aleix",
	//}
	os.Chdir("/data/work/go/GoDevEach/ansible/ansibleplaybook-extravars-file")

	var res *results.AnsiblePlaybookJSONResults

	buff := new(bytes.Buffer)
	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		//Inventory: "127.0.0.1,",
		Forks:     "10",
		Inventory: "hosts", // hostfile path
		ExtraVarsFile: []string{
			"@vars-file1.yml",
			"@vars-file2.yml",
		},
		ExtraVars: map[string]interface{}{options.AnsibleHostKeyCheckingEnv: false},
	}

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "smart",
	}

	playbook_fusion := &playbook.AnsiblePlaybookCmd{
		Playbooks: []string{"site.yml"},
		//ConnectionOptions: ansiblePlaybookConnectionOptions,
		Exec: execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
		Options:           ansiblePlaybookOptions,
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		StdoutCallback:    "json",
	}

	options.AnsibleAvoidHostKeyChecking()
	options.AnsibleForceColor()
	err := playbook_fusion.Run(context.TODO())
	if err != nil {
		fmt.Println(err) //可能 是其中 一台主机的错误
	}

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		panic(err)
	}
	//panic(len(res.Plays))
	for _, play := range res.Plays {
		for _, task := range play.Tasks {
			//fmt.Println(task.Task.Name)
			for host, content := range task.Hosts {
				//fmt.Println(host)
				//fmt.Printf("%+v\n", content)
				//fmt.Printf("host: %s, %+v,%v\n", host, content.Stdout, content.Stderr)
				fmt.Printf("Task [%s] Host: [%s] Msg: [%v] Action: [%s] failed [%t] with condition [%t]. Executed cmd: %v,Stdout: %v, Stderr: %v\n",
					task.Task.Name, host, content.Msg, content.Action, content.Failed, content.FailedWhenResult, content.Cmd, content.Stdout, content.Stderr)
				fmt.Println("========================================================")
			}
		}
	}
}
