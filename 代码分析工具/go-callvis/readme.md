安装

```
ubuntu
apt-get install graphviz

go get -u github.com/ofabry/go-callvis
# or
git clone https://github.com/ofabry/go-callvis.git
cd go-callvis && make install
```

官方网站

[https://github.com/ofabry/go-callvis#readme](https://github.com/ofabry/go-callvis#readme)

使用


```
cd /home/go/GoDevEach/works/datacenter_go
go-callvis . | dot -Tpng -o syncthing.png
```

