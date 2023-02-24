package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
)

func main() {
	os.Chdir("/data/work/go/GoDevEach/ansible/ansible-loopResult")
	var err error
	var res *results.AnsiblePlaybookJSONResults

	buff := new(bytes.Buffer)
	timeBuff := new(bytes.Buffer)

	ansiblePlaybookConnectionOptions := &options.AnsibleConnectionOptions{
		Connection: "local",
	}

	ansiblePlaybookOptions := &playbook.AnsiblePlaybookOptions{
		Inventory: "127.0.0.1,",
	}
	execute.NewDefaultExecute()
	playbooksList := []string{"site1.yml"}
	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks: playbooksList,
		Exec: execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		StdoutCallback:    "json",
	}

	err = playbook.Run(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
	}

	res, err = results.ParseJSONResultsStream(io.Reader(buff))
	if err != nil {
		panic(err)
	}

	for _, play := range res.Plays {
		for _, task := range play.Tasks {
			for host, content := range task.Hosts {
				fmt.Println(host)

				if task.Task.Name == "skipping-task" {
					fmt.Printf("Task [%s] skipped [%t] with skip reason [%s]\n",
						task.Task.Name, content.Skipped, content.SkipReason)
				} else {
					fmt.Printf("Task [%s] failed [%t] with condition [%t]. Executed cmd: %v\n",
						task.Task.Name, content.Failed, content.FailedWhenResult, content.Cmd)
				}

			}
		}
	}

	fmt.Println(timeBuff.String())
}