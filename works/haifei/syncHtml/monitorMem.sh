#!/bin/bash
d=`date "+%F %T"`
m=`ps aux|grep _build|grep -v grep `
str5="date: ${d} ; Mem: ${m}"
echo $str5 >> /root/2
