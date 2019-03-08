## 描述 ##
SSDB是一个高性能的支持丰富数据结构的NoSQL数据库。系统使用SSDB主要用于保存各种配置项，如后端服务器信息、集群信息、路由信息等等。SSDB可以作单机部署，也可以作集群部署。作集群部署时，配置为主主模式。

**SSDB特性**

- 替代 Redis 数据库, Redis 的 100 倍容量。
- LevelDB 网络支持, 使用 C/C++ 开发。
- Redis API 兼容, 支持 Redis 客户端。
- 适合存储集合数据, 如 list, hash, zset...。
- 客户端 API 支持的语言包括: C++, PHP, Python, Java, Go。
- 持久化的队列服务。
- 主从复制, 负载均衡。

## 安装 ##
直接解压ssdb.tar.gz文件，解压后目录为ssdb。使用命令赋予主目录ssdb读写权限。

## 启动 ##
在主目录中运行  
		
	./start.sh

	成功输出：
	ssdb-server 1.9.4
	Copyright (c) 2012-2015 ssdb.io

## 停止 ##
在主目录中运行

	./stop.sh
	
	成功输出：
	ssdb-server 1.9.4
	Copyright (c) 2012-2015 ssdb.io

## 配置文件(ssdb.conf)说明 ##

    work_dir = ./var  			//工作目录
    pidfile = ./var/ssdb.pid  	//存放ssdb进程pid的文件目录
    server:
        ip: 0.0.0.0   			//本机ip
	    port: 18888  			//本机服务端口

    replication:  
	    binlog: yes  
	    sync_speed: -1  
	    slaveof:  				//主从模式，屏蔽后表示不适用主从或主主模式
		    #id: ssdb_1  		
		    #type: mirror  
		    #host: 192.168.10.31  
		    #port: 18888  

    logger:  
	    level: error  
	    output: stdout  
	    rotate:  
		    size: 1000000000  

    leveldb: 
	    cache_size: 500  
	    write_buffer_size: 64  
	    compaction_speed: 1000  
	    compression: yes  


**必须修改项**

	server:
        ip: 0.0.0.0   			//机器有多个网卡时，最好配置一个固定ip
	    port: 18888  			//本机服务端口


## 将SSDB配置为主主模式 ##
假设有2台机器安装SSDB  
**server1(192.168.10.1:18888)**配置文件如下：
 	
	server:
        ip: 0.0.0.0   			
	    port: 18888  			
    replication:  
	    binlog: yes  
	    sync_speed: -1  
	    slaveof:  				
		    id: ssdb_1  		
		    type: mirror  
		    host: 192.168.10.2  //配置为对方ip
		    port: 18888  

**server2(192.168.10.2:18888)**配置文件如下：
 	
	server:
        ip: 0.0.0.0   			
	    port: 18888  			
    replication:  
	    binlog: yes  
	    sync_speed: -1  
	    slaveof:  				
		    id: ssdb_2  		
		    type: mirror  
		    host: 192.168.10.1  //配置为对方ip
		    port: 18888  

**备注：配置修改后需要重启SSDB**