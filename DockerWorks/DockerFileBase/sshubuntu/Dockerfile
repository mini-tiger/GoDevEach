FROM harbor.dev.21vianet.com/taojun/ubuntu:20.04

ADD sources.list /etc/apt
RUN  apt update && apt-get install -y openssh-server net-tools && sed -i 's/#Port 22/Port 22/' /etc/ssh/sshd_config && mkdir /run/sshd \
    &&  sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && apt-get clean \
    && mkdir -p /root/.ssh && ssh-keygen -q -t rsa -b 2048 -f /root/.ssh/id_rsa -N '' \
    && cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys && chmod 400 /root/.ssh/authorized_keys \
    && cat /root/.ssh/id_rsa

SHELL ["/bin/bash", "-c"]
RUN echo "root:123456" | chpasswd

CMD ["/usr/sbin/sshd", "-D"]

#docker run -itd -p 2222:22 {imagename}