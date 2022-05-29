## Nano DB
一款以golang编写的轻量、非关系型、基于系统文件系统的kv数据库。  
说白了就是把数据按文件存在本地，再提供http接口来访问，因此可以适用于分布式服务。  
在日常使用的服务器上，不会如同常见数据库：速度会随着数据量增大而“显著”减慢。   

## 特点
- 轻量：即使包含数十万索引，树莓派上也能流畅运行
- 高速：微秒级查询
- RESTful接口：不熟悉SQL语句，没问题
- 缓存：查询结果缓存，提高查询效率
- 权限管理：ACL，每个用户权限分离

## 使用
### CLI总览
```sh
Usage of ./nano-db:
  -a string
        specific the addr to listen (default "0.0.0.0:3777")
  -d string
        update acl rules with -d <dbname>
  -l int
        set the max length of cache (default 100)
  -s string
        set salt for cookie
  -u string
        generate the cookie with -n <username>
```

### 获取cookie
`./nano-db -c {userName}`  
为你的用户生成cookie  
cookie会被打印到控制台，请在后继操作时，在headers内附带此cookie

### 启动数据库
`./nano-db`
可以使用`-a`参数指定监听地址，默认为`0.0.0.0:3777`  
使用`-l`参数指定缓存的最大长度，默认为100

### 数据库操作
#### 查看数据库是否存活
`HEAD /`  
唯一不需要附带cookie的接口，可用于客户端检查数据库是否存活  

#### 查看总状态
`GET /`
会输出有多少数据库、COL、内存缓存项及获取时间

#### 初始化
`./nano-db -u {userName} -d {dbName}`
需要先初始化数据库，才能进行后继操作  
第一个初始化{DB}的用户将会成为该{DB}的唯一管理员  

可以打开`.acl/acl.json`文件进行手动修改  
例如：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel"]}]}
```
你想给用户`novel`添加访问数据库`test`的权限，可以如下修改：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel","test"]}]}
```

⚠️**注意**，如果当前数据库正在运行，acl更改将在一分钟内应用。

#### 获取DB内所有Col
`GET /{DB}`

#### 删除数据库
`DELETE /{DB}`

#### 获取Col内所有ID
`GET /{DB}/{COL}`

#### 删除某Col
`DELETE /{DB}/{COL}`

#### 获取
`GET /{DB}/{COL}/{ID}`

#### 插入/更新
`POST /{DB}/{COL}/{ID}`
需要在body附带需要写入的数据

#### 删除
`DELETE /{DB}/{COL}/{ID}`


## 注意⚠️
`{DB}`,`{COL}`,`{ID}` 不能包含字符 `/` ` ` `\\` `..`，并且他们的长度都不能超过37.

