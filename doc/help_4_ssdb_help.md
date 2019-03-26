# 命令行工具 ssdb-cli #
SSDB 的命令行工具 `ssdb-cli` 对于 SSDB 的管理非常有用,你可以用它来执行所有的命令,监控服务的状态,清除整个数据库,等等。

## 连接到 SSDB 服务器 ##

    [root@xq25 ssdb]# ./ssdb-cli -p 16888
	ssdb (cli) - ssdb command line tool.
	Copyright (c) 2012-2016 ssdb.io

	'h' or 'help' for help, 'q' to quit.

	ssdb-server 1.9.4

	ssdb 127.0.0.1:16888> 

# 常用命令 #
## info ##
返回服务器的信息
    
	ssdb 127.0.0.1:16888> info
	version
        1.9.4
	links
        1
	total_calls
        15220433
	dbsize
        0
	binlogs
            capacity : 20000000
            min_seq  : 0
            max_seq  : 335
	replication
        client 192.168.10.90:12373
            type     : mirror
            status   : SYNC
            last_seq : 335
	replication
        slaveof 192.168.10.90:16888
            id         : ssdb_1
            type       : mirror
            status     : OUT_OF_SYNC
            last_seq   : 862505
            copy_count : 72
            sync_count : 0
	serv_key_range
            kv  : "" - ""
            hash: "" - ""
            zset: "" - ""
            list: "" - ""
	data_key_range
            kv  : "" - ""
            hash: "HashCluster" - "HashServer"
            zset: "" - ""
            list: "" - ""
	leveldb.stats
                                       Compactions
        Level  Files Size(MB) Time(sec) Read(MB) Write(MB)

## Key-Value ##
**set key value** 设置指定 key 的值内容
    
	ssdb 127.0.0.1:16888> set key1 value1
	ok
	(0.001 sec)

**get key** 获取指定 key 的值内容
	
	ssdb 127.0.0.1:16888> get key1
	value1
	(0.001 sec)	

**del key** 删除指定的 key
	
	ssdb 127.0.0.1:16888> del key1
	ok
	(0.001 sec)

## Hashmap ##
**hset name key value** 设置 hashmap 中指定 key 对应的值内容

	ssdb 127.0.0.1:16888> hset HashName1 key1 value1
	ok
	(0.001 sec)

**hget name key** 获取 hashmap 中指定 key 的值内容

	ssdb 127.0.0.1:16888> hget HashName1 key1
	value1
	(0.001 sec)

**hsize name** 返回 hashmap 中的元素个数
	
	ssdb 127.0.0.1:16888> hsize HashName1 
	1
	(0.000 sec)

**hgetall name** 返回整个 hashmap

	ssdb 127.0.0.1:16888> hgetall HashName1
	key             value
	-------------------------
  	key1           : value1
	1 result(s) (0.000 sec)
	(0.000 sec)

**hdel name key** 删除 hashmap 中的指定 key
	
	ssdb 127.0.0.1:16888> hdel HashName1 key1
	1
	(0.002 sec)

**hclear name** 删除 hashmap 中的所有 key
 
	ssdb 127.0.0.1:16888> hclear HashName1
	0
	(0.003 sec)