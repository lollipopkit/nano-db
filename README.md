## Nano DB
一款以golang编写的轻量kv数据库。

## 特点
- RESTful接口：不熟悉SQL语句，没问题
- 缓存：查询结果缓存，提高查询效率
- 权限管理：ACL，每个用户权限分离
- 轻量：甚至可以在Raspberry Pi Zero上无压力运行
- 高速：微秒级查询

## 使用
### CLI总览
```sh
Usage of nano-db:
  -c string
        generate the cookie with -c <username>
  -l int
        set the max length of cache (default 100)
  -s string
        set salt for cookie
  -u string
        specific the addr to listen (default "0.0.0.0:3777")
```
### 更改salt
两种方法：
- 修改`consts/app.go`内`CookieSalt`的值，需要固定的值，随意填写。  
- 使用`-s`参数在运行时指定。例如：`./nano-db -s "1234567890"`


### 获取cookie
`./nano-db -c {userName}`  
为你的用户生成cookie  
在执行该步骤前请确认是否完成了上一步（修改salt）
然后cookie会被打印到控制台，请在后继操作时，在headers内附带此cookie

### 数据库
#### 查看数据库是否存活
`HEAD /`  
唯一不需要鉴权的接口  

#### 查看总状态
`GET /`
会输出有多少数据库、COL、缓存项及获取时间

#### 初始化
`HEAD /{DB}`
需要先初始化数据库，才能进行后继操作  
第一个初始化{DB}的用户将会成为该数据库的唯一管理员  

如果你想手动管理权限，可以打开`.acl/acl.json`文件进行手动修改  
例如：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel"]}]}
```
你想给用户`novel`添加访问数据库`test`的权限，可以如下修改：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel","test"]}]}
```

#### 获取DB内所有Collection
`GET /{DB}`

#### 删除数据库
`DELETE /{DB}`

#### 获取指定Collection内所有ID
`GET /{DB}/{COL}`

#### 删除某Col
`DELETE /{DB}/{COL}`

#### 是否存在
`HEAD /{DB}/{COL}/{ID}`

#### 获取
`GET /{DB}/{COL}/{ID}`

#### 插入/更新
`POST /{DB}/{COL}/{ID}`
需要在body附带需要写入的数据

#### 删除
`DELETE /{DB}/{COL}/{ID}`


## 注意⚠️
`{DB}`,`{TABLE}`,`{ID}` 不能包含字符 `/`,` `,`\\`,`..`，并且他们的长度都不能超过37.

