package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Monitor struct {
	Cmd    string `json:"cmd"`
	Return int    `json:"return"`
}

type InstallAll struct {
	Interval    int                 `json:"interval"`
	StopErr     bool                `json:"stopErr"`
	RunDir      string              `json:"run_dir"`
	InstallStep map[int]interface{} `json:"install_step"`
}

type Install struct {
	InstallAll InstallAll `json:"install_all"`
	Monitor    Monitor    `json:"monitor"`
	StartCmd   string     `json:"start_cmd"`
	StopCmd    string     `json:"stop_cmd"`
}

type InstallResp struct {
	Success bool
	error   string
	Step    string
	cmd     string
}

func mapKeySort(m map[int]interface{}) (keys []int) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return
}

func UnJson(jb []byte) (err error, dl *Install) {
	err = json.Unmarshal(jb, &dl) //解析
	if err != nil {
		log.Println("解析json失败")
		return err, dl
	}
	return nil, dl
}

func main() {
	f, err := ioutil.ReadFile("/home/go/src/godev/mymodels/运行系统命令/install.json")
	if err != nil {
		fmt.Println(err)
	}
	err, dataJson := UnJson(f) // 解析Json

	if err != nil {
		fmt.Println(1, err)
	}
	fmt.Println(dataJson)
	fmt.Println(dataJson.Monitor)

	cinstall := make(chan struct{}, 0)

	for _, step := range []string{"installStep"} {
		resp, stop := StepRun(step, dataJson, cinstall)
		if stop {
			fmt.Println(resp)
			break
		} else {
			fmt.Println(11111, resp)
		}
	}
	//resp:=StepRun("installStep", dataJson, cinstall)
	//
	//fmt.Println(resp)
	//resp=StepRun("startStep", dataJson,cinstall)
	////
	////<-cinstall
	////<-cinstall
	//
	//fmt.Println(resp)
	//cstart := make(chan struct{},0)
	//
	//RunCommand()

}

func StepRun(step string, dl *Install, c chan struct{}) (resp InstallResp, stop bool) {
	//defer func() {
	//	c <- struct{}{}
	//}()
	var mapInstStep interface{}
	switch step {
	case "installStep":
		mapInstStep = dl.InstallAll.InstallStep
		mapInstStep = mapInstStep.(map[int]interface{})
	}

	mapInst := mapInstStep.(map[int]interface{})
	slInstStep := mapKeySort(mapInst)

	for _, v := range slInstStep {
		cmd, _ := mapInst[v]
		stderr, stdout := RunCommand(cmd.(string))
		if stderr != "" && dl.InstallAll.StopErr {
			resp.error = strings.TrimRight(stderr, "\n")
			resp.cmd = cmd.(string)
			resp.Step = step
			resp.Success = false
			stop = true
			break
		}
		log.Printf("begin run Step:%s, index: %d,stdout:%s,stderr:%s\n", step, v, stdout, stderr)
		time.Sleep(time.Duration(dl.InstallAll.Interval) * time.Second)
	}
	return
}

func RunCommand(command string) (string, string) {
	cmd := exec.Command("/bin/bash", "-c", command)

	var out bytes.Buffer
	cmd.Stdout = &out

	var e bytes.Buffer
	cmd.Stderr = &e
	//cmd.Start()
	//buf, _ := cmd.CombinedOutput()
	cmd.Run()
	if e.Len() != 0 && out.Len() == 0 {
		return e.String(), out.String()
	}
	return "", out.String()
}
