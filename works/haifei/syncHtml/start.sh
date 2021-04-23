#!/bin/bash
if [ `whoami` != "root" ];then
	echo -e "\033[41;37m 当前不是root用户,启动失败! \033[0m"
	exit 1
fi
HomeDir="/home/go/GoDevEach/works/haifei/syncHtml"
binName="synchtml"
pushd $HomeDir

    #启动方法    
    start(){  
     #now=`date "+%Y%m%d%H%M%S"`  
     nohup ${HomeDir}/${binName} > /dev/null 2>&1 &
    }  
    #停止方法  
    stop(){  
     ps -ef|grep ${binName}|grep -v grep|awk '{print $2}'|while read pid  
     do  
        kill -9 $pid  
     done  
    }  
      
    case "$1" in  
    start)  
    start  
    ;;  
    stop)  
    stop  
    ;;    
    restart)  
    stop  
    start  
    ;;  
    *)  
    printf 'Usage: %s {start|stop|restart}\n' "$prog"  
    exit 1  
    ;;  
    esac
