## Nano DB
一款以golang编写的轻量、非关系型、基于文件系统的kv数据库。  

白话文：将数据按文件储存，再提供http接口来访问，因此可以适用于分布式服务（一台数据库服务器，多个后端服务器）。  


## 特点
- 轻量：即使包含数十万索引，树莓派上也能流畅运行
- 高速：微秒级查询
- RESTful接口：无需SQL语句（目前：意味着没有where、order by等）
- 缓存：查询结果缓存，提高查询效率
- 权限管理：ACL，每个用户权限分离
- SDK：目前支持 [go](https://git.lolli.tech/lollipopkit/nano-db-sdk-go)

## 使用
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

#### 启动数据库
`./nano-db`
可以使用`-a`参数指定监听地址，默认为`0.0.0.0:3777`  
使用`-l`参数指定缓存的最大长度，默认为100

#### 获取cookie
`./nano-db -c {userName}`  

为你的用户生成cookie，cookie会被打印到控制台  
⚠️ **请在使用http接口时，在headers内附带此cookie。或以此cookie使用sdk**

#### 添加权限
`./nano-db -d {dbName} -u {userName}`   

指定用户成为指定数据库的唯一管理员  

可以打开`.sct/acl.json`（如文件不存在，需要先启动数据库一次）文件进行手动修改，例如：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel"]}]}
```
你想给用户`novel`添加访问数据库`test`的权限，可以如下修改：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel","test"]}]}
```

⚠️**注意**，如果当前数据库正在运行，acl更改将在一分钟内应用。


### 数据库操作
操作数据库可以选择：
- SDK（[go](https://git.lolli.tech/lollipopkit/nano-db-sdk-go)，其他sdk待开发）
- HTTP接口

接下来是http接口的使用，sdk文档请前往sdk查看。


方法|接口|功能|额外说明
---|---|---|---
HEAD|`/`|查看数据库是否存活|唯一不需要附带cookie的接口，可用于客户端检查数据库是否存活
GET|`/`|查看总状态|会输出有多少数据库、COL、内存缓存项及获取时间
GET|`/{DB}`|获取DB内所有Col|会返回所有col的名称，并非db内所有col的数据
DELETE|`/{DB}`|删除数据库|不会删除对该数据库的权限
GET|`/{DB}/{COL}`|获取Col内所有ID|获取col下所有id的名称，并非col下所有数据
DELETE|`/{DB}/{COL}`|删除某Col|并且删除col下所有ID
GET|`/{DB}/{COL}/{ID}`|获取|不存在则会返回错误
POST|`/{DB}/{COL}/{ID}`|插入/更新|需要在body附带需要写入的数据
DELETE|`/{DB}/{COL}/{ID}`|删除|如果路径不存在则会返回错误
⚠️**注意**：`{DB}`,`{COL}`,`{ID}` 不能包含字符 `/` ` ` `\` `..`，并且他们的长度都不能超过37。

