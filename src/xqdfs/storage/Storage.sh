#!/bin/bash
mode=$1  # start or stop
dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ulimit -c 0

Checkkeeper()
{
    if [ -f $dir/Storage.pid ];then
	keeppid=`cat $dir/Storage.pid`
	kill -0 $keeppid >/dev/null 2>&1
	keepstatus1="$?"
        #echo "keepstatus1 is $keepstatus1"
	keepstatus2=`ps -ef | grep "/Storage start" | grep -v grep | wc -l` 
	#echo "keepstatus2 is $keepstatus2"
	if [ $keepstatus1 -ne 0 ] || [ $keepstatus2 -eq 2 ];then	
	    return 0
	else
	    return 1
	fi
    else
	    return 0	
    fi
}

CheckappStorage()
{
    if [ -f $dir/storage.pid ];then
	apppid=`cat $dir/storage.pid`
	kill -0 $apppid >/dev/null 2>&1
	appstatus1="$?"
	appstatus2=`ps -ef | grep "storage" | grep -v grep | wc -l`
    	if [ $appstatus1 -ne 0 ] || [ $appstatus2 -eq 0 ]; then
       	    return 0
    	else
            return 1
    	fi
    else
       return 0
    fi
}

case $mode in
    'start')
       Checkkeeper
       Checkkeeper_RET=$?
       if [ $Checkkeeper_RET -ne 0 ]; then
	  echo -e "\033[1;31mStorage has already running,if you want to restart,please stop first!\033[0m"
	  echo -e "Storage's pid is \033[1;31m$keeppid\033[0m"
	  exit
       fi
       cat /dev/null >$dir/starttimes.txt
	   i=1
       while true
       do
		
	  CheckappStorage
	  if [ $? -eq 0 ]; then
			cd $dir
			./storage &
			echo $! >storage.pid
			echo "storage" >>$dir/starttimes.txt
          fi
          sleep 5
       done &
       keeperpid=$(jobs -p)
       echo "$keeperpid" >$dir/Storage.pid
    ;;
    'stop')
    # stop Service
    	if [ ! -f $dir/Storage.pid ];then
	    echo "Storage is not running!"
	    exit
    	fi

	kill -9 `cat $dir/Storage.pid`
    rm -f $dir/Storage.pid
	
	kill -3 `cat $dir/storage.pid`
	rm -f $dir/storage.pid
        
    ;;
    *)
	echo "Usage: $0 {start|stop}"
        exit 1
    ;;
esac
exit 0
