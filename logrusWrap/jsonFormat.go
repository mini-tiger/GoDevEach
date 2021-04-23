package main

import (
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{ // xxx 自定义默认列的列名
		log.FieldKeyTime:  "@timestamp",
		log.FieldKeyLevel: "@level",
		log.FieldKeyMsg:   "@message"},})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	logger = log.WithFields(log.Fields{"request_id": "123444", "user_ip": "127.0.0.1"})




}

func main() {
	logger.Info("hello, logrus....")
	logger.Info("hello, logrus1....")
	log.WithFields(log.Fields{
	 "animal": "walrus",
	 "size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(log.Fields{
	 "omg":    true,
	 "number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(log.Fields{
	 "omg":    true,
	 "number": 100,
	}).Fatal("The ice breaks!")
}