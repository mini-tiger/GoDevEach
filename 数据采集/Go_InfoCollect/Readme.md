
```shell
cd /data/work/go/GoDevEach/数据采集/

docker run --rm \
-v $(pwd)/Go_InfoCollect:/data/main \
harbor.dev.21vianet.com/taojun/golang:static1.17 \
bash -c "go build -mod=vendor -ldflags=' -linkmode external -extldflags "-static"' -a -o main"


```

