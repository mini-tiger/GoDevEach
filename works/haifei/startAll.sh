#!/bin/bash
basedir="/home/go/GoDevEach/works/haifei"
source /etc/profile

lsof -i:1521
if [ `echo $?` -gt 0 ]
then
        sleep 5
fi

#redis:
redis-server /etc/redis.conf
#backend:
pushd /home/work/icework-hf-proj/target
sh start.sh stop
sh start.sh start
#front:
#systemctl restart httpd
pushd /home/work/haifei
sh startClient.sh

pushd ${basedir}/MonitorCommcell
sh start.sh stop
sh start.sh start
popd

pushd ${basedir}/syncHtmlYWReport
sh start.sh stop
sh start.sh start
popd

