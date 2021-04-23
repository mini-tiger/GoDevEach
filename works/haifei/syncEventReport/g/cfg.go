package g

import (
	"log"
	"github.com/toolkits/file"
	"sync"
)

var (

	lock       = new(sync.RWMutex)

)

func ParseConfig(cfg string) string{
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)

	}

	//var c GlobalConfig


	lock.Lock()
	defer lock.Unlock()

	log.Println("read config file:", cfg, "successfully")
	return configContent
	//WLog(fmt.Sprintf("read config file: %s successfully",cfg))
}
