package main

import (
	logdiy "gitee.com/taojun319/tjtools/logDiyNew"
)

func main() {
	//var Log *Log1

	Log := logdiy.InitLog1("C:\\work\\go-dev\\AutoNomy\\nxlogNew\\1.log", 3, true, "WARNING")
	//Log = Logger()
	Log.Printf("%d\n", 1111)
	Log.Printf("%d\n", 1111)
	Log.Error("%d\n", 1111)
	Log.Debug("%d\n", 1111)
}
