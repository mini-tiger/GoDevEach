package g

import (
	"encoding/json"
	"gitee.com/taojun319/tjtools/file"
	"log"
	"sync"
)

//var PassMap []string
//var HostData []string
var configLock = new(sync.RWMutex)
var config *ConfigJson

func Config() *ConfigJson {

	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

type ConfigJson struct {
	Hosts     []interface{} `json:"hosts"`
	PasswdMap []interface{} `json:"passwd_map"`
}

func ParseConfig(cfg string) {

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c ConfigJson
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
	log.Println("config hosts:", config.Hosts)
	log.Println("config passwd:", config.PasswdMap)
}
