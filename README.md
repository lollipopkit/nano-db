## Nano DB
一款以golang编写的轻量kv数据库。

## 特点
- RESTful接口
- 缓存
- 权限管理
- 轻量：甚至可以在Raspberry Pi Zero上运行

## 使用
### 更改salt
修改`consts/app.go`内`CookieSalt`的值，需要固定的值，随意填写。
### 获取cookie
为你的用户生成cookie  
`./nano-db -c {userName}`  
然后cookie会被打印到控制台，后继需要鉴权的操作，都需要在headers内附带此cookie

### 数据库
#### 查看数据库是否存活
`HEAD /`

#### 查看总状态
`GET /`

#### 初始化
需要先初始化数据库，才能进行后继操作  
第一个访问{DB}的用户将会成为该数据库的唯一管理员  
`GET /{DB}`

#### 获取
`GET /{DB}/{COL}/{ID}`

#### 插入/更新
`POST /{DB}/{COL}/{ID}`
需要在body附带需要写入的数据（仅支持json）

#### 删除
`DELETE /{DB}/{COL}/{ID}`


## 注意⚠️
`{DB}`,`{TABLE}`,`{ID}` 不能包含字符 `/`,` `,`\\`,`..`，并且他们的长度都不能超过20.

