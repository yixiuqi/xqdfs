## 安装SSDB
解压ssdb.tar.gz，修改配置文件ssdb.conf

	server:
		ip: 0.0.0.0
		port: 18888 改为打算使用的端口


启动SSDB 
	
	./start.sh

启动后会打印 

	ssdb-server 1.9.4
	Copyright (c) 2012-2015 ssdb.io


停止SSDB 
	
	./stop.sh


## 启动Storage
解压Storage.tar.gz，修改配置文件xqdfs_storage.toml

	[Log]
	Level = "error"            日志等级trace，debug，info，error
	
	[Server]
	Id = 1                     当前节点编号，编号必须唯一   
	Desc = "test[25]"          当前节点描述
	Host = "192.168.10.25"     本机IP地址
	Port = 10086               本地使用端口
	
	[Dir]
	Path = ["volume_data"]
	Capacity = [128]           本机用于存储的空间(单位G)，默认128G
	
	[Configure]
	Param = "192.168.10.25:18888"    SSDB地址

启动 
	
	./Storage.sh start  

停止 
	
	./Storage.sh stop

## 启动Master
解压Master.tar.gz，修改配置文件xqdfs_master.toml

	[Log]
	Level = "error"                 日志等级trace，debug，info，error
	
	[Server]
	Host = "192.168.10.25"          本机IP地址                                    
	Port = 10087                    本地使用端口
	
	[Configure]
	Param = "192.168.10.25:18888"   SSDB地址

启动

	 ./Master.sh start**  

停止
	
	 ./Master.sh stop**