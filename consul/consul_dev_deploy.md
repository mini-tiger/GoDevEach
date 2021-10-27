
### deploy
```
$ wget https://releases.hashicorp.com/consul/1.7.3/consul_1.7.3_linux_amd64.zip
$ mv consul /usr/local/bin/
$ consul -v
Consul v1.7.3
```

### 发布服务

先启动 192.168.1.153 上一个测试Web服务
```
$ vim /etc/consul.d/web.json #探测发现服务的配置文件【注册服务】
{

  "services": [
    {

      "id": "hello1",
      "name": "hello",
      "tags": [
        "primary"
      ],
      "address": "192.168.43.111",
      "port": 8080,
      "checks": [
        {

        "http": "http://192.168.43.111:8080/jenkins/",
        "tls_skip_verify": false,
        "method": "Get",
        "interval": "10s",
        "timeout": "1s"
        }
      ]
    },
   {

      "id": "hello2",
      "name": "hello",
      "tags": [
        "primary"
      ],
      "address": "192.168.43.111",
      "port": 8080,
      "checks": [
        {

        "http": "http://192.168.43.111:8080/jenkin/",
        "tls_skip_verify": false,
        "method": "Get",
        "interval": "10s",
        "timeout": "1s"
        }
      ]
    }

  ]
}


$ mkdir /etc/consul/data #存储数据目录
$ consul agent -dev -ui -client 0.0.0.0 -config-dir /etc/consul.d -data-dir=/etc/consul/data
或者
$ consul agent -dev -ui -client 0.0.0.0 -config-file /etc/consul.d/web.json -data-dir=/etc/consul/data
```

### 访问
```
http://192.168.43.111:8500/
```

### 获取服务

```
192.168.43.111 是运行consul 服务器
hello 是web.json 定义的

方式一
http://192.168.43.111:8500/v1/catalog/service/hello

方式二
 dig @192.168.43.111 -p 8600 hello.service.consul
```
