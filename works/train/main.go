package main

import (
	"flag"
	"fmt"
	"godev/works/train/g"
	"godev/works/train/db"
	"os"
	"os/signal"
	"syscall"
	"godev/works/train/cron"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")

	flag.Parse()

	g.ParseConfig(*cfg)

	g.InitRedis()
	g.InitLog1()

	db.Init()

	cron.TrainCrond()

	//cache.Init()

	//go cache.DeleteStaleAgents()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		//db.DB.Close()
		os.Exit(0)
	}()

	select {}
}
