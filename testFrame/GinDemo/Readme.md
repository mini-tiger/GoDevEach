```shell

go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega/...

cd GinDemo && mkdir test
cd test
创建入口文件  
ginkgo bootstrap

创建单元测试文件
ginkgo generate controllers // 随便起名

cd test
go test or  ginkgo -p -v

```

```shell
https://ke-chain.github.io/ginkgodoc/#%E5%B9%B6%E8%A1%8C-specs
https://github.com/onsi/gomega
https://github.com/onsi/ginkgo

```