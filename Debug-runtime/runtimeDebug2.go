package main

import (
  "fmt"
  "log"
  "runtime"
  "runtime/debug"
  "time"

  "github.com/google/gops/agent"
)

func main() {
  if err := agent.Listen(agent.Options{
    Addr: "0.0.0.0:8848",
    // ConfigDir:       "/home/centos/gopsconfig", // 最好使用默认
    ShutdownCleanup: true}); err != nil {
    log.Fatal(err)
  }

  fmt.Println(debug.SetGCPercent(1))

  // 1
  var dic = make([]byte,100,100)
  runtime.SetFinalizer(&dic, func(dic *[]byte) {
    fmt.Println("内存回收1")
  })

  // 立即回收
  runtime.GC()

  // 2
  var s = make([]byte,100,100)
  runtime.SetFinalizer(&s, func(dic *[]byte) {
    fmt.Println("内存回收2")
  })
  //runtime.GC()
  // 3
  d := make([]byte,300,300)
  for index,_ := range d {
    d[index] = 'a'
  }
  fmt.Println(len(d))
  runtime.SetFinalizer(&d, func(dic *[]byte) {
    fmt.Println("内存回收3")
  })

  time.Sleep(100*time.Second)
}
