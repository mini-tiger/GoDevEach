package main

import (
  "runtime"
  "fmt"
    "time"
  "runtime/debug"
)


type A int

func main() {

  //debug.SetGCPercent(1)
  var a *A=new(A)

  runtime.SetFinalizer(a, func(d *A) {
    fmt.Println("内存回收a",d)
  })


  var dic = new(map[string]string)
  runtime.SetFinalizer(dic, func(d *map[string]string) { //dict 在回收前的触发
    fmt.Println("内存回收dic",d)
  })

  runtime.GC()
  debug.FreeOSMemory()


  time.Sleep(time.Second)
}
