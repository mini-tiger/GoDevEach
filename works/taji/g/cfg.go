package g

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/taojun319/tjtools/file"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	//ReadFile    = "10.155.2.4_yuebao_English.html"
	//htmlfileReg = "C:\\work\\go-dev\\src\\godev\\works\\haifei\\*.html"
	//htmlBakDir  = "C:\\work\\go-dev\\src\\godev\\works\\haifei\\bak\\"
	//timeinter   = 10

	FieldsLen = 21 // 列一共有几列,使用时 要 减一

	//ConfigJson = "D:\\work\\project-dev\\GoDevEach\\works\\haifei\\syncHtml\\synchtml.json"

)

var cfg Config
var configLock = new(sync.RWMutex)

type Config struct {
	Debug      bool   `json:"debug"`
	Stdout     bool   `json:"stdout"`
	RedisDB    int    `json:"redisDB"`
	LogMaxDays int    `json:"logMaxDays"`
	RedisPwd   string `json:"redisPwd"`
	RedisAddr  string `json:"redisAddr"`
	Logfile    string `json:"logfile"`
	LogLevel   string `json:"loglevel"`
	SubChan    string `json:"SubChan"`
	MongoUri   string `json:"mongoUri"`
}

func ParseConfig(cfg string) string {
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

func CheckConfig(fp string) (e error, conf string) {
	// 兼容开发与生产环境

	if file.IsExist(fp) {
		return nil, fp
	} else {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		if file.IsExist(filepath.Join(dir, fp)) {
			return nil, filepath.Join(dir, fp)
		} else {
			return errors.New(fmt.Sprintf("confile :%s Not Found", fp)), ""
		}
	}

}

func readconfig(cfgfile string) {
	cfgstr := ParseConfig(cfgfile)
	err := json.Unmarshal([]byte(cfgstr), &cfg)
	if err != nil {
		log.Fatalln("parse config file fail:", err)
	}
}

func LoadConfig(cfgPath string) {

	e, confile := CheckConfig(cfgPath) // xxx 更改配置文件使用上面const
	if e == nil {
		readconfig(confile)
		log.Printf("config file success:%+v\n", cfg)
	} else {
		log.Fatalln("config file fail:", e)
	}
}

func GetConfig() *Config {

	configLock.RLock()
	defer configLock.RUnlock()
	return &cfg
}
