package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/apenella/go-ansible/pkg/adhoc"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/apenella/go-ansible/pkg/execute"
	"github.com/apenella/go-ansible/pkg/options"
	"github.com/apenella/go-ansible/pkg/playbook"
	"github.com/apenella/go-ansible/pkg/stdoutcallback/results"
)

var parentCtx context.Context = context.Background()

type MissionEntry struct {
	*ConnParams
	Success   bool      `json:"Success"`
	Err       error     `json:"Err"`
	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`
}

type ConnParams struct {
	ctx      context.Context
	IP       string `json:"IP"`
	Pass     string `json:"pass"`
	User     string `json:"User"`
	KeyFile  string `json:"KeyFile"`
	Port     string `json:"port"`
	AuthType string `json:"auth_type"` // pass or keyfile

}

func main() {
	//byte1, err := ioutil.ReadFile("25keyfile")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(byte1))
	c := &MissionEntry{
		ConnParams: &ConnParams{
			IP:       "172.22.50.25",
			Pass:     "",
			User:     "root",
			KeyFile:  "25keyfile", // 文件路径
			Port:     "32468",
			AuthType: "keyfile",
		},
		Success: false,
		Err:     nil,
	}
	c.RunMission()
	cc := &MissionEntry{
		ConnParams: &ConnParams{
			IP:       "172.22.50.22",
			Pass:     "123456",
			User:     "root",
			KeyFile:  "",
			Port:     "32139",
			AuthType: "pass",
		},
		Success: false,
		Err:     nil,
	}
	cc.RunMission()
	time.Sleep(15 * time.Second)
	c.MissionResultToDB()
	cc.MissionResultToDB()
}

func (c *MissionEntry) RunMission() {
	ctx, _ := context.WithTimeout(parentCtx, time.Duration(10*time.Second))
	c.ctx = ctx
	go c.SingleMission()
	//time.Sleep(10 * time.Second)
	//cancel()
}

func (c *MissionEntry) MissionResultToDB() {
	if c.Success {
		log.Printf("ip:%s,port:%s starttime:%v, endtime:%v\n", c.IP, c.Port, c.StartTime, c.EndTime)
	} else {
		log.Printf("ip:%s,port:%s starttime:%v, endtime:%v err:%v\n", c.IP, c.Port, c.StartTime, c.EndTime, c.Err)
	}

}

func (c *MissionEntry) MissionResultErr(err string) {
	c.Err = errors.New(err)
}

func (c *MissionEntry) SingleMission() {
	c.StartTime = time.Now()
	defer func() {
		c.EndTime = time.Now()
	}()

	defer func() {
		if err := recover(); err != nil {
			c.Success = false
			c.MissionResultErr(fmt.Sprintf("run panic err: %v", err))
			return
		}
	}()
	os.Chdir("/data/work/go/GoDevEach/ansible/ansible-golang-plugin")
	var err error
	var res *results.AnsiblePlaybookJSONResults

	buff := new(bytes.Buffer)
	//timeBuff := new(bytes.Buffer)

	options.AnsibleAvoidHostKeyChecking() //  "ANSIBLE_HOST_KEY_CHECKING": "false"

	optVars := map[string]interface{}{"ansible_ssh_pass": c.Pass, "ansible_ssh_port": c.Port, "ansible_ssh_user": c.User}

	var ansiblePlaybookOptions *playbook.AnsiblePlaybookOptions
	var ansiblePlaybookConnectionOptions *options.AnsibleConnectionOptions

	if strings.Contains(c.AuthType, "key") {
		// key
		// ssh -i 25keyfile root@172.22.50.25 -p 32468 -o StrictHostKeyChecking=no
		ansiblePlaybookOptions = &playbook.AnsiblePlaybookOptions{
			Inventory: fmt.Sprintf("%s,", c.IP),
			ExtraVars: optVars,
		}
		ansiblePlaybookConnectionOptions = &options.AnsibleConnectionOptions{
			PrivateKey: c.KeyFile,
			Connection: "smart",
			//SSHCommonArgs: "StrictHostKeyChecking=no",
		}

	} else { // pass

		// 密码
		ansiblePlaybookOptions = &playbook.AnsiblePlaybookOptions{
			//Inventory: "172.22.50.21, 172.22.50.25,",
			Inventory: fmt.Sprintf("%s,", c.IP),
			ExtraVars: optVars,
			ExtraVarsFile: []string{
				"@vars-file1.yml",
				"@vars-file2.yml",
			},
		}

		ansiblePlaybookConnectionOptions = &options.AnsibleConnectionOptions{
			//PrivateKey: "25keyfile",
			Connection: "smart",
			//SSHCommonArgs: "StrictHostKeyChecking=no",
		}

	}

	//execute.NewDefaultExecute()
	playbooksList := []string{"site.yml"}
	playbook := &playbook.AnsiblePlaybookCmd{
		Playbooks: playbooksList,
		Exec: execute.NewDefaultExecute(
			execute.WithWrite(io.Writer(buff)),
		),
		ConnectionOptions: ansiblePlaybookConnectionOptions,
		Options:           ansiblePlaybookOptions,
		StdoutCallback:    "json",
	}
	for {
		select {
		case <-c.ctx.Done():
			c.Success = false
			c.MissionResultErr(fmt.Sprintf("run timeout err: %v", ""))
			return
		default:
			err = c.adhocPing(ansiblePlaybookConnectionOptions, optVars)
			if err != nil {
				c.Success = false
				//fmt.Println("2222222222222222222", err.Error())
				c.MissionResultErr(fmt.Sprintf("adhocPing err: %v", err.Error()))
				return
			}
			err = playbook.Run(c.ctx)
			if err != nil {
				c.Success = false
				c.MissionResultErr(fmt.Sprintf("playbook err: %v", err.Error()))
				return
			}

			res, err = results.ParseJSONResultsStream(io.Reader(buff))
			if err != nil {

				c.Success = false
				c.MissionResultErr(fmt.Sprintf("result ParseJson err: %v", err.Error()))
				return
			}

			for _, play := range res.Plays {
				for _, task := range play.Tasks {
					for host, content := range task.Hosts {
						fmt.Println(host)
						//if task.Task.Name == "skipping-task" {
						//	fmt.Printf("Task [%s] skipped [%t] with skip reason [%s]\n",
						//		task.Task.Name, content.Skipped, content.SkipReason)
						//} else {
						//	fmt.Printf("Task [%s] failed [%t] with condition [%t]. Executed cmd: %v\n",
						//		task.Task.Name, content.Failed, content.FailedWhenResult, content.Cmd)
						//}
						fmt.Printf("task stdout:%v stderr:%v\n", content.Stdout, content.Stderr)

					}
				}
			}
			c.Success = true
			return

		}

	}

	//fmt.Println(timeBuff.String())
}

func (c *MissionEntry) adhocPing(connOpt *options.AnsibleConnectionOptions, optionVars map[string]interface{}) error {

	//ansibleConnectionOptions := &options.AnsibleConnectionOptions{
	//	Connection: "smart",
	//}

	ansibleAdhocOptions := &adhoc.AnsibleAdhocOptions{
		Inventory:  fmt.Sprintf("%s,", c.IP),
		ModuleName: "ping",
		ExtraVars:  optionVars,
	}
	//buff := new(bytes.Buffer)
	//errbuff := new(bytes.Buffer)
	adhoc1 := &adhoc.AnsibleAdhocCmd{
		Pattern:           "all",
		Options:           ansibleAdhocOptions,
		ConnectionOptions: connOpt,
		StdoutCallback:    "json",
		//Exec: execute.NewDefaultExecute(
		//	execute.WithWrite(io.Writer(buff)),
		//	execute.WithWriteError(io.Writer(errbuff)),
		//),
	}

	log.Printf("ip:%s adhoc Command:\n %v\n ", c.IP, adhoc1.String())
	//fmt.Println("!!!!!!!!!!11", buff.String(), errbuff.String())
	return adhoc1.Run(c.ctx)

}
