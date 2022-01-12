#! /bin/sh  
HomeDir="/home/go/GoDevEach/works/haifei/MonitorCommcell"
binName="monitorCommCell"
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
    #重载配置文件
    reload(){
     ps -ef|grep ${binName}|grep -v grep|awk '{print $2}'|while read pid
     do
        kill -10 $pid
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
    reload)
    reload
    ;;
    *)  
    printf 'Usage: %s {start|stop|restart|reload}\n' "$prog"
    exit 1  
    ;;  
    esac
