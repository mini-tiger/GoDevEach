1. 测试 1
```shell
go build main
第一次启动
/bin/bash start.sh

kill -1 pid
ps aux|grep -E "main|nohup"
```


2. 测试 2 (不需要shell 启动)
```shell
go build -o main1 main1.go 
第一次启动
./main1 > /tmp/22 2>&1 &
kill -1 pid
ps aux|grep main1
```