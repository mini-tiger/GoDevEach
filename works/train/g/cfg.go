package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
	"github.com/go-xorm/xorm"
)

type redisinfo struct {
	Dsn      string `json:"dsn"`
	Password string `json:"password"`
	PoolSize int    `json:"poolSize"`
}

type Train struct {
	LoopInterval   int `json:"loopinterval"`   // train 表到达traininterval 间隔后,没有找到新车进入，循环扫描库的间隔
	TrainInterval  int `json:"traininterval"`  // 根据上一条 train pass_time 加上这个时间，下次获取train的时间
	MaxVehicleTime int `json:"maxvehicletime"` // 根据train_serial，每次vehicle 等待的时间
}
type UdpInfo struct {
	UdpAddr string `json:"udpAddr"`
}
type GlobalConfig struct {
	Debug     bool       `json:"debug"`
	OracleDsn string     `json:"oracledsn"`
	Redis     *redisinfo `json:"redis"`

	Logfile    string   `json:"logfile"`
	LogMaxDays int      `json:"logMaxDays"`
	Daemon     bool     `json:"daemon"`
	Train      *Train   `json:"train"`
	UdpInfo    *UdpInfo `json:"udpInfo"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
	lock       = new(sync.RWMutex)
	Engine     *xorm.Engine //定义引擎全局变量
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")
}
