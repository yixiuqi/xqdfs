## XQDFS介绍 ##
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;XQDFS是一个高性能分布式文件系统。 它的主要功能包括：文件上传、文件下载、文件清理，以及高容量和负载平衡。主要解决了海量数据存储问题，特别适合以中小文件为载体的在线服务。  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;XQDFS系统包括：调度服务器(Master)、存储服务器(Storage)、客户端(Client)、Web管理界面。

* **Master**：调度服务器，管理所有存储服务器，负载均衡作用。  
* **Storage**：存储服务器，主要提供容量和备份服务。  
* **Client**：客户端，上传下载数据的服务器，也就是我们自己的项目所部署在的服务器。

## 版本 ##
**V0.0.1.20190311 alpha**

## 部署模式 ##
![典型部署方式****](image/logic.png "")