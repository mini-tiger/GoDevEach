
```shell
docker build -f Dockerfile -t aa:aa .

docker run --rm -it  -v $(pwd)/main:/data/main aa:aa bash
# -a 重新编译 不用缓存 
 go build -ldflags=' -linkmode external -extldflags "-static"' -a -o main
```
