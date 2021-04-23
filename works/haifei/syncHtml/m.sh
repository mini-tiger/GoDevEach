#!/bin/bash
prog_name="synchtml"
prog_mem=$(ps aux|grep ${prog_name}|grep -v grep|awk '{print $11}')
time=$(date "+%Y-%m-%d %H:%M:%S")

for i in ${prog_mem}
do
	m=$(ps aux|grep ${i}|grep -v grep|awk '{print $6}')
	echo "time":${time}"; Mem":${m}"; progname:"${i} >> /root/a11123sync
done
