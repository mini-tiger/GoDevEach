#FROM harbor.dev.21vianet.com/taojun/ubuntu:20.04
FROM docker.io/library/ubuntu:20.04
ADD sources.list /etc/apt
# https://ftp.gnu.org/gnu/glibc/
ADD glibc-2.32.tar.gz /opt/
# https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-2.1.1-linux-x86_64.tar.bz2
ADD phantomjs-2.1.1-linux-x86_64/bin/phantomjs /usr/bin/

RUN apt update && apt-get install -y build-essential python3.10 gawk bison python3-pip && \
    ln -s /usr/bin/python3.10 /usr/bin/python

RUN apt remove -y libc6-dev && mkdir -p /opt/glibc && cd /opt/glibc && \
    ../glibc-2.32/configure --prefix=/opt/glibc --disable-werror  --without-selinux --enable-add-ons  --enable-static-nss  --enable-kernel=2.6.10 && \
    make -j 8 && make install && apt-get install -y libc6-dev && \
    pip config set global.index-url http://mirrors.cloud.tencent.com/pypi/simple && pip config set global.trusted-host mirrors.cloud.tencent.com && \
    pip install snapshot-phantomjs==0.0.3 && pip3 install Pillow==10.0.1 &&  pip install pyecharts==2.0.4  && apt-get install -y xfonts-wqy


ENV LIBRARY_PATH=/opt/glibc/lib:$LIBRARY_PATH
ENV OPENSSL_CONF=/dev/null
WORKDIR /data/main
