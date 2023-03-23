#!/bin/bash
set -e
pushd $(pwd) > /dev/null
    cd ${GOPATH}/GoDevEach/打包编译多目录父目录子目录
    DIRS=$(find * -maxdepth 0 -type d|grep -v package)

    for tmp in $DIRS;do
        FILES=$(find $tmp -name 'Makefile')
        echo $FILES
        for tmp_file in $FILES;do
            target_makefile_path=$(pwd)/$tmp_file
            if [ -f $target_makefile_path ];then
                pushd $(pwd) > /dev/null
                    cd $(dirname $target_makefile_path)
		    echo "enter directory: " $(pwd)
                    if [ "$1" = "debug" ];then
                        export ISDEBUG=true
                    fi

                    make -f Makefile
                    if [ $? -ne 0 ];then
                        exit
                    fi
                popd > /dev/null
            fi
        done
    done
popd > /dev/null

