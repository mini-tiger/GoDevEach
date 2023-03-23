1. first.go
```shell
执行文件 外部调用者 不承载业务代码 
```

2. demo1.go
```shell

go build demo1.go
./demo1
tail -f demo1.log
ps aux|grep demo1
```

3. deamon/demo.go 第三方库
```shell
go build demo.go
./demo
ps aux|grep demo

```