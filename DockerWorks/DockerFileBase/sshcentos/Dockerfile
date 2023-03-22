FROM harbor.dev.21vianet.com/taojun/centos:7.4.1708-debug
RUN sed -e 's|^mirrorlist=|#mirrorlist=|g' \
             -e 's|^#baseurl=http://mirror.centos.org|baseurl=https://mirrors.tuna.tsinghua.edu.cn|g' \
             -i.bak \
             /etc/yum.repos.d/CentOS-*.repo \
    && yum install -y net-tools  openssl openssh-server \
    && ssh-keygen -q -t rsa -b 2048 -f /etc/ssh/ssh_host_rsa_key -N '' && ssh-keygen -q -t ecdsa -f /etc/ssh/ssh_host_ecdsa_key -N '' \
    && ssh-keygen -t dsa -f /etc/ssh/ssh_host_ed25519_key -N '' &&  sed -i 's/#Port 22/Port 22/' /etc/ssh/sshd_config \
    && sed -i 's/#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config && echo "123456" | passwd --stdin root >& /dev/null \
    && yum clean all && mkdir -p /root/.ssh && cat /etc/ssh/ssh_host_rsa_key.pub >> /root/.ssh/authorized_keys && chmod 400 /root/.ssh/authorized_keys \
    && cat /etc/ssh/ssh_host_rsa_key

CMD ["/usr/sbin/sshd", "-D"]

#docker run -itd -p 2222:22 {imagename}