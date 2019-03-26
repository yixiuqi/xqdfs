## 导入系统预置参数 ##
系统所有配置都是保存在SSDB中，SSDB初始安装后数据库是空的，需要手动导入系统默认的配置项。

## 启动SSDB客户端 ##
在SSDB主目录中运行

	./ssdb-cli -p 18888		
	
	输出：
	ssdb (cli) - ssdb command line tool.
	Copyright (c) 2012-2016 ssdb.io

	'h' or 'help' for help, 'q' to quit.

	ssdb-server 1.9.4

	ssdb 127.0.0.1:18888>

18888是端口号，这个根据服务端配置的端口号为准

## 导入默认配置 ##
将**backup.ssdb**文件(系统预置配置文件)放到主目录中，然后运行

	ssdb 127.0.0.1:18888> import backup.ssdb
	
	输出：
	 5%
	10%
	21%
	44%
	86%
	96%
	100%
	done.


