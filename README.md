## 简介
IM 即时聊天demo

技术: Gin + Websocket + Redis + Mysql + MongoDB

数据库:
- MySQL: 存储用户基本信息
- MongoDB: 存放用户聊天信息
- Redis: 存储 session 和 用户发送消息的数量(骚扰消息拦截)

Websocket:
- 服务端与客户端通信

## 项目功能
- 用户管理
    - 注册
    - 登录
    - 修改密码
    - 登出
- 双人聊天
    - 在线、不在线应答
    - 获取历史聊天记录
    - 获取未读聊天记录
    - 删除聊天记录
    - 修改聊天记录为已读
    - 骚扰消息拦截
- 在线群聊

## 项目结构
- api: API实现
- cache: Redis缓存, session和发送信息的数量
- conf: 项目配置, mysql(存储用户信息)、monogo db(存储聊天信息)
- model: mysql、monogo db 的数据库操作
- pkg/status: 状态码及其说明
- pkg/utils: 工具函数, 包括加密、解密和自定义错误等
- router: 路由转发
- serializer: http消息序列化
- service: 服务实现, 包括用户服务、websocket通信实现

## 配置文件
conf/config.ini
```ini 
# debug开发模式,release生产模式
[service]
AppMode = debug
# 运行端口号 8081 端口
HttpPort = 8081

[mysql]
Db = mysql
# mysql的ip地址
DbHost = "127.0.0.1"
# mysql的端口号,默认3306
DbPort = 3306
# mysql user
DbUser = test_root
# mysql password
DbPassWord = 123456
# 数据库名字
DbName = chat_demo

[redis]
# redis 名字
RedisDb = redis
# redis 地址
RedisAddr = "127.0.0.1:6379"
# redis 密码
RedisPw = ""
# redis 数据库名
RedisDbName = 2

[MongoDB]
# 数据库名称
MongoDBName = chat_demo
# 地址
MongoDBAddr = localhost
# 用户
MongoDBUser = test_root
# 密码
MongoDBPwd = 123456
# 端口
MongoDBPort = 27017
```

## 项目运行
```bash
go mod tidy
go run main.go
```

## 接口
### http
测试服务是否可以连接

Get: http://127.0.0.1:8081/ping

聊天服务

Get: http://127.0.0.1:8081/chat/test?uid=user1&to_uid=user2

注册服务

Post: http://127.0.0.1:8081/user/register

![register.png](images%2Fregister.png)

登录服务

Post: http://127.0.0.1:8081/user/login

![login.png](images%2Flogin.png)

设置密码

Post: http://127.0.0.1:8081/user/set_password

![set_password.png](images%2Fset_password.png)

退出登录

Post: http://127.0.0.1:8081/user/logout

![logout.png](images%2Flogout.png)

### websocket
```bash
ws://127.0.0.1:8081/chat/test?uid=1&to_uid=2
```
单聊: content数据格式: 消息id + \n + 消息内容
```json
{
    "from":"1",
    "to":"2",
    "type": 1,
    "op_type":1,
    "content":"消息id\nhello"
}
```
群聊
```json
{
    "from":"1",
    "to":"",
    "type": 1,
    "op_type":6,
    "content":"hello everyone"
}
```
获取与对应用户的聊天记录: content数据格式: 跳过数据-获取数据数量
```json
{
    "from":"1",
    "to":"2",
    "type": 1,
    "op_type":2,
    "content":"3-5"
}
```
获取与对应用户的所有未读消息
```json
{
    "from":"1",
    "to":"2",
    "type": 1,
    "op_type":3,
    "content":""
}
```
将某一个消息 ID 修改为已读: content: 消息id
```json
{
    "from":"1",
    "to":"",
    "type": 1,
    "op_type":4,
    "content":"消息id"
}
```
删除对应的消息: content: 消息id
```json
{
    "from":"1",
    "to":"",
    "type": 1,
    "op_type":5,
    "content":"消息id"
}
```

[参考](https://github.com/CocaineCong/gin-chat-demo)
