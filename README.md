# agent



代理主要为BOX系统私钥APP提供服务，转发voucher、company、appserver、私钥APP等的信息。

代理的环境依赖与编译环境：

| 名称         | 说明       | 备注 |
| ------------ | ---------- | ---- |
| 使用环境依赖 | mysql      |      |
| 编译环境     | golang,git |      |



### 下载&编译

```shell
#下载
mkdir -p $GOPATH/src/github.com/boxproject 
cd $GOPATH/src/github.com/boxproject
git clone https://github.com/boxproject/agent.git

#编译
go build
```



### 运行

编译完成后会产生agent文件，agent运行需要的文件结构

```shell
$GOPATH/src/github.com/boxproject/agent/
├── agent							#可执行文件
├──conf
│  └──  app.conf					#HTTPS服务配置文件
├── certs							#GRPC证书存放文件夹
├── config.json						#系统配置、appserver配置文件
├── leveldb							#leveldb数据文件夹
├── log								#日志存放文件夹
├── log.xml							#日志参数配置文件
├── scripts							#证书生成脚本文件,手动生成证书需要
│   ├── client-cert.sh
│   └── server-cert.sh
└── ssl								#HTTPS服务证书存放文件夹
```



创建必要文件夹，创建配置文件

```shell
cd $GOPATH/src/github.com/boxproject/agent/
mkdir log leveldb ssl conf 
cp config.json.example  config.json
cp log.xml.example  log.xml
cp app.conf.example conf/app.conf
```



##### 文件配置

config.json配置文件，参考模版config.json.example，修改“/opt/box/agent”为自己当前的工程目录

```htmp
#数据库相关配置
"url":"[agent数据库用户名]:[agent数据库用户密码]@tcp([数据库IP地址]:3306)/[agent数据库database]?charset=utf8&loc=Asia%2FShanghai",
...
#leveldb相关配置
"level_db_path":"/opt/box/agent/leveldb",
#GRPC证书相关配置
"server_cert":"/opt/box/agent/certs/server.pem",
"server_key":"/opt/box/agent/certs/server.key",
"client_cert":"/opt/box/agent/certs/client.pem",
#appserver服务对应的IP及端口,，deposit_url、withdraw_url等参数
"manager_ip_port":"app服务的IP:5001",  
"deposit_url":"http://app服务的IP:5001/api/v1/capital/deposit",
...
```

HTTPS服务配置文件app.conf，参考模版app.conf.example，修改“/opt/box/agent”为自己当前的工程目录

```shell
EnableHTTP = false
EnableHTTPS = true
EnableHttpTLS = true
HTTPSPort = 19092
HTTPSCertFile = "/opt/box/agent/ssl/server.pem"
HTTPSKeyFile = "/opt/box/agent/ssl/server.key"
```

agent日志配置文件log.xml，参考模版log.xml.example，主要更改level、property两个地方，修改“/opt/box/agent”为自己当前的工程目录，更加详细的使用方法参考[log4go](https://github.com/alecthomas/log4go)，

```html
...
<filter enabled="true">
    <tag>file</tag>
    <type>file</type>
    <level>DEBUG</level>
	<property name="filename">/opt/box/agent/log/agent.log</property>
...
```



##### 证书的产生

```shell
#grpc证书
#client-cert后期给voucher和companion使用
cd $GOPATH/src/github.com/boxproject/agent/
./scripts/client-cert.sh ./certs
./scripts/server-cert.sh ./certs

#https 证书
$./server-cert.sh ./ssl
```



##### 数据库数据表的构建

下载sql脚本文件：【[db_table.sql](https://github.com/boxproject/agent/blob/master/db_table.sql) 】

```shell
#1，登录mysql
$mysql -u root -p 
#2，切换db
$mysql>USE boxagent; 
#3，更新数据库
$mysql>source  [path]/db_table.sql

```



### 使用方法

```shell
#启动
cd /opt/box/agent 
./agent start 
#停止
./agent stop
```

