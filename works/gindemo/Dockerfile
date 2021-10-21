#xxxx 设置基础镜像,Harbor服务器pull镜像
FROM harbor.od.com/library/glibc:bash
#创建app目录,保存我们的代码
RUN mkdir -p /usr/src/GinDemo
#设置工作目录
WORKDIR /usr/src/GinDemo

#复制所有文件到 工作目录。 pwd当前目录是
COPY ./main /usr/src/GinDemo/


#暴露container的端口
EXPOSE 8001
ENV mysqluser root
#运行命令
CMD ["sh","-c","./main >server.log 2>&1"]
