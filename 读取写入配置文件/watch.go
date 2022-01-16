package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

/**
 * @Author: Tao Jun
 * @Description: main
 * @File:  watch
 * @Version: 1.0.0
 * @Date: 2021/11/15 下午6:30
 */

const defaultConfigFile = "config.yaml"

var ViperCfg Config

func main() {
	select {}
}

func init() {
	v := viper.New()
	v.SetConfigFile("/home/go/GoDevEach/读取配置文件/config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&ViperCfg); err != nil {
			fmt.Println(err)
		}
		log.Println(ViperCfg)
	})
	if err := v.Unmarshal(&ViperCfg); err != nil {
		fmt.Println(err)
	}
	log.Println(ViperCfg)
}

type Config struct {
	Mysql      Mysql      `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis      Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
	System     System     `mapstructure:"system" json:"system" yaml:"system"`
	JWT        JWT        `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Log        Log        `mapstructure:"log" json:"log" yaml:"log"`
	Prometheus Prometheus `mapstructure:"prometheus" json:"prometheus" yaml:"prometheus"`
	Consul     Consul     `json:"consul" yaml:"consul"`
	Public     Public     `mapstructure:"public" json:"public" yaml:"public"`
}

type Prometheus struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
}

type System struct {
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	Env           string `mapstructure:"env" json:"env" yaml:"env"`
	Addr          string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
}

type JWT struct {
	SigningKey string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`
}

type Casbin struct {
	ModelPath string `mapstructure:"model-path" json:"modelPath" yaml:"model-path"`
}

type Mysql struct {
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

type Log struct {
	Prefix  string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	LogFile bool   `mapstructure:"log-file" json:"logFile" yaml:"log-file"`
	Stdout  string `mapstructure:"stdout" json:"stdout" yaml:"stdout"`
	File    string `mapstructure:"file" json:"file" yaml:"file"`
}

type Sqlite struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Path     string `mapstructure:"path" json:"path" yaml:"path"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
	LogMode  bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
}

type Consul struct {
	Addr string `json:"addr" yaml:"addr"`
}

type Public struct {
	PublicSendEmailAddress string `mapstructure:"publicSendEmailAddress" json:"publicSendEmailAddress" yaml:"publicSendEmailAddress"`
}
