#!/bin/sh
# stworker 服务进程管理
# author zhangpf@chuchujie.com
# $Id: server.sh 2017-02-13

SERVICE_FLAG=vcutters
SYSTEM_TIME=`date '+%Y-%m-%d %T'`
cd $(dirname $0) && {
    SHELL_PATH=`pwd`
    cd - &> /dev/null
}

GOIL_HOME=${SHELL_PATH}
SHELL_LOG_PATH=$GOIL_HOME/logs/shell
SHELL_PIDFILE=$GOIL_HOME/logs/shell/shell.pid
LOG=$GOIL_HOME/logs/shell/shell.log
PID_FILE=$GOIL_HOME/logs/shell/vcutter_pid.log

if [ ! -d "$SHELL_LOG_PATH" ]; then
    mkdir "$SHELL_LOG_PATH"
fi
if [ ! -f "$PID_FILE" ]; then
    touch "$PID_FILE"
fi
# if [ ! -f "$SHELL_PIDFILE" ]; then
#     touch "$SHELL_PIDFILE"
# fi

function check_shell {
    SHELL_PID=`cat $SHELL_PIDFILE 2> /dev/null`
    if [ $? -eq 1 ]; then
        echo 'no pid file'
    else
        count=`ps aux | grep -w $SHELL_PID | grep 'server.sh' | grep -v grep | wc -l`
        if [ $count -gt 0 ]; then
            echo "shell is running  pid : $SHELL_PID" >> $LOG
            exit
        fi
    fi
}

date >> $LOG
check_shell

flock -n $SHELL_PIDFILE -c "echo $$ > $SHELL_PIDFILE; sleep 0.5"

if [ $? -eq 1 ]; then 
    SHELL_PID=`cat $SHELL_PIDFILE 2> /dev/null`
    echo "$SHELL_PID is running" >> $LOG
    exit
fi
echo "shell pid is $$" >> $LOG

function start {
    cd $GOIL_HOME && {
        pid=`cat $PID_FILE`
        count=`ps aux | grep -e "${SERVICE_FLAG}$" | grep -v grep | wc -l`
        if [ $count -eq 1 ]; then
            echo "vcutter is running, pid is $pid"
        else
            nohup ${SHELL_PATH}/${SERVICE_FLAG} > /dev/null 2>&1
            if [ $? -eq 0 ]; then
                g_pid=`ps aux | grep -e "${SERVICE_FLAG}$" | grep -v grep | awk '{print $2}'`
                echo $g_pid>${PID_FILE}
                echo "vcutter start success, pid is `cat $PID_FILE`"
            else
                echo "vcutter start fail"
            fi
        fi
        exit 0
    }
}

function stop {
    cd $GOIL_HOME && {
        pid=`cat $PID_FILE`
        count=`ps aux | grep -e "${SERVICE_FLAG}$" | grep -v grep | wc -l`
        echo "stop vcutter , the number of running process is $count" >> $LOG
        if [ $count -ge 1 ]; then
            `ps -ef | grep -e "${SERVICE_FLAG}$" | grep -v "grep" | awk '{print $2}' | xargs sudo kill -9`
            echo "vcutter is stoping" >> $LOG
        else
            echo "stop fail, vcutter is not running" >> $LOG
        fi
        exit 0
    }
}

function reload {
    cd $GOIL_HOME && {
        pid=`ps aux | grep -e "${SERVICE_FLAG}$" | grep -v grep | awk '{print $2}'`
        if [[ $? != 0  || -z $pid ]]; then
            echo "vcutter pid is null" >> $LOG
            start
            exit 0
        fi
        `kill -12 $pid`
        if [ $? -ge 1 ]; then
            echo "vcutter reload fail" >> $LOG
        else
            for p in `ps aux | grep ${SERVICE_FLAG} | grep -v grep | awk '{print $2}'`;do
                if [ $pid -ne $p ];then
                    echo $p > ${PID_FILE}
                fi
            done
            echo "vcutter reload successfully" >> $LOG
        fi
    }
}

case $1 in

    start)
    start
    exit;;

    stop)
    stop 
    exit;;

    restart)
    stop
    start
    exit;;

    reload)
    reload
    exit;;

    *)
    echo "invalid arg: $1" >> $LOG;;
esac
