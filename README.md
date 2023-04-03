## Nano DB
一款以golang编写的轻量、非关系型、基于文件系统的数据库。  

白话文：将数据按文件储存，再提供http接口来访问。  

## 🔖 特点
- 无需SQL语句：使用gjson与正则搜索匹配数据
- 轻量：即使包含数十万索引，树莓派上也能流畅运行
- 高速：微秒级查询
- RESTful接口：HTTP协议，方便使用
- 缓存：查询结果缓存，提高查询效率
- 权限管理：ACL，每个用户权限分离
- SDK：目前支持 [go](https://github.com/lollipopkit/nano-db-sdk-go)

## 📖 使用
```sh
Usage of ./nano-db:
  -d string
        update acl rules with -d <dbname>
  -u string
        generate the cookie with -u <username>
```

#### 启动数据库
`./nano-db`
可以编辑 `.cfg/app.json` 修改配置

#### 获取cookie
`./nano-db -u {userName}`  

为你的用户生成cookie，cookie会被打印到控制台  

⚠️ **请在使用http接口时，在headers内附带此cookie。或以此cookie使用sdk**

#### 添加权限
`./nano-db -d {dbName} -u {userName}`   

指定用户成为指定数据库的唯一管理员  

可以打开 `.cfg/acl.json`（如文件不存在，需要先启动数据库一次或手动创建）文件进行手动修改，例如：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel"]}]}
```
如果想给用户 `novel` 添加 `test` 数据库的权限，可以如下修改：
```json
{"ver":1,"rules":[{"user":"novel","db":["novel","test"]}]}
```

⚠️**注意**，如果当前数据库正在运行，acl更改将在一分钟内应用。


### 🔨 数据库操作
操作数据库可以选择：
- SDK（[go](https://github.com/lollipopkit/nano-db-sdk-go)，其他sdk待开发）
- HTTP接口

接下来是http接口的使用，sdk文档请前往sdk查看。


方法|接口|功能|额外说明
---|---|---|---
HEAD|`/`|查看数据库是否存活|唯一不需要附带cookie的接口，可用于客户端检查数据库是否存活
GET|`/`|查看总状态|会输出有多少数据库、DIR、内存缓存项及获取时间
GET|`/{DB}`|获取DB内所有DIR|会返回所有DIR的名称，并非DB内所有DIR的数据
DELETE|`/{DB}`|删除数据库|不会删除对该数据库的权限
POST|`/{DB}`|搜索DB下所有文件|返回包含 `gjson.Get(FILE,p).Exists()` 为真的文件内容。如果正则 `v` 不为空，则会剔除 `gjson.Result.Raw` 不匹配的。body结构：`{"path":"","regex":""}`
GET|`/{DB}/{DIR}`|获取DIR内所有FILE|获取DIR下所有文件的名称，并非DIR下所有数据
DELETE|`/{DB}/{DIR}`|删除某DIR|并且删除DIR下所有FILE
POST|`/{DB}/{DIR}`|搜索DIR下所有文件|body结构：`{"path":"","regex":""}`。如果正则 `v` 为空，返回包含 `gjson.GetBytes(FILE,path).Exists()` (`FILE` 为 `DIR` 下文件的内容)，为真的文件内容。如果正则 `v` 不为空，则会剔除 `regexp.MatchString(regex, gjson.Result.Raw)` 不匹配的。
GET|`/{DB}/{DIR}/{FILE}`|获取|不存在则会返回错误
POST|`/{DB}/{DIR}/{FILE}`|插入/更新|需要在body附带需要写入的数据
DELETE|`/{DB}/{DIR}/{FILE}`|删除|如果路径不存在则会返回错误

⚠️**注意**：`{DB}`,`{DIR}`,`{FILE}` 不能包含字符 `/` ` ` `\` `..`，并且他们的长度都不能超过37。

建议规范：`novel/chapter/1.json` `xapp/user/xxx.json` `secret/key/xxx.json`

## 🔒 安全
请妥善保管你的cookie，不要将其发送给他人。  
如果发现cookie被盗，请更改 `.cfg/app.json` 内的 `Salt`。  
随后重新生成cookie，并使用新的cookie访问数据库。

## 🔑 License
`LGPL LollipopKit 2022`