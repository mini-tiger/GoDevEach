FROM golang:1.17 AS build
COPY . /go/src/matchAllRoutes
WORKDIR /go/src/matchAllRoutes
ENV GOPROXY https://goproxy.cn,direct
RUN make matchAllRoutes.docker

#xxxx 设置基础镜像
#FROM gongt/glibc:bash
FROM scratch
#FROM alpine:3.13.6
# xxx size  glibc > alpine > scratch
# xxx glibc alpine 可以增加包
#创建app目录,保存我们的代码
#RUN mkdir -p /usr/src/go

#复制所有文件到 工作目录。 pwd当前目录是
COPY --from=build /go/src/matchAllRoutes/matchAllRoutes.docker /usr/src/

#设置工作目录
WORKDIR /usr/src/
#COPY ./config.json /usr/src/datacenter_go/
# xxx alpine 时区
#RUN apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata

# xxx scratch 时区
# ./Dockerfile_tz

#暴露container的端口
EXPOSE 8081
#ENV mysqluser root
#运行命令
#CMD ["sh","-c","./main >server.log 2>&1"]

# 日志打印在 标准输出
# xxx if from glibc:bash
#CMD ["bash","-c","./matchAllRoutes.docker"]
# xxx if from scratch
ENTRYPOINT ["./matchAllRoutes.docker"]


