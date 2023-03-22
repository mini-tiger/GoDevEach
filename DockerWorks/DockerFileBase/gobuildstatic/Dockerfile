#FROM harbor.dev.21vianet.com/taojun/ubuntu:20.04
FROM docker.io/library/ubuntu:20.04
ADD sources.list /etc/apt

ADD go1.17.13.linux-amd64.tar.gz /usr/local/
ADD glibc-2.32.tar.gz /opt/


RUN apt update && apt-get install -y build-essential python3 gawk bison && \
    apt remove -y libc6-dev && \
    mkdir -p /opt/glibc && cd /opt/glibc && \
    ../glibc-2.32/configure --prefix=/opt/glibc --disable-werror  --without-selinux --enable-add-ons  --enable-static-nss  --enable-kernel=2.6.10 && \
    make -j 8 && make install  && \
    apt-get install -y libc6-dev

ENV GOROOT=/usr/local/go \
    GOPROXY=https://goproxy.cn,direct \
    GOPATH=/data \
    PATH=$PATH:/usr/local/go/bin



#RUN apt remove -y libc6-dev && \
#    mkdir -p /opt/glibc && cd /opt/glibc && \
#     ../glibc-2.32/configure --prefix=/opt/glibc --disable-werror  --without-selinux --enable-add-ons  --enable-static-nss  --enable-kernel=2.6.10 && \
#    make -j 8 && make install  \
#    && apt-get install -y libc6-dev

ENV LIBRARY_PATH=/opt/glibc/lib:$LIBRARY_PATH

WORKDIR /data/main
