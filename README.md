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

### 启动数据库
`nano-db`
可以编辑 `.cfg/app.json` 修改配置

### 添加权限

#### 设置权限
`nano-db -d {dbName}`  
示例：
```
➜  nano-db git:(main) ✗ nano-db -d novel
[INF] generated token: ijlV5aKzMja0MgTkpd0Q8J6zuegtwzVQzEd8A
[SUC] acl update rule: success
```
可以用生成的 `token` 访问名为 `novel` 的数据库

⚠️ **注意**：
- 请在使用http接口时，在 `headers` 内的 `NanoDB` 键附带此 `token`
- 使用 sdk，需要用到此 `token` 

#### 添加权限

##### 命令行添加
`nano-db -d {dbName} -t {token}`   

示例：
```
➜  nano-db git:(main) ✗ nano-db -d novel -t ijlV5aKzMja0MgTkpd0Q8J6zuegtwzVQzEd8A
[SUC] acl update rule: success
```
然后就可以用该 `token` 访问名叫 `novel` 的数据库了

##### 手动添加
如果不愿意将 `token` 暴露至 shell，可以打开 `.cfg/acl.json`（如文件不存在，需要先启动数据库一次或手动创建）文件进行手动修改，例如：
```json
{
      "ver":1,
      "rules":[
            {
                  "token":"token1",
                  "dbs":[
                        "novel"
                  ]
            }
      ]
}
```

如果想给 `token1` 添加 `test` 数据库的权限，可以如下修改：
```json
{
      "ver":1,
      "rules":[
            {
                  "token": "token1",
                  "dbs":[
                        "novel",
                        "test"
                  ]
            }
      ]
}
```

⚠️**注意**：
如果当前数据库正在运行，acl更改将在一分钟内应用。


### 🔨 数据库操作
接下来是http接口的使用，sdk文档请前往sdk查看。

方法|接口|功能|额外说明
---|---|---|---
HEAD|`/`|查看数据库是否存活|唯一不需要附带cookie的接口，可用于客户端检查数据库是否存活
GET|`/{DB}`|获取DB内所有DIR|会返回所有DIR的名称，并非DB内所有DIR的数据
DELETE|`/{DB}`|删除数据库|不会删除对该数据库的权限
GET|`/{DB}/{DIR}`|获取DIR内所有FILE|获取DIR下所有文件的名称，并非DIR下所有数据
DELETE|`/{DB}/{DIR}`|删除某DIR|并且删除DIR下所有FILE
GET|`/{DB}/{DIR}/{FILE}`|获取|不存在则会返回错误
POST|`/{DB}/{DIR}/{FILE}`|更新|需要在body附带需要写入的数据
DELETE|`/{DB}/{DIR}/{FILE}`|删除|如果路径不存在则会返回错误

⚠️**注意**：`{DB}`,`{DIR}`,`{FILE}` 不能包含除 `.`(不能在开头结尾) `0-9` `A-Z` `a-z` 以外的字符，并且他们的长度都不能超过37。

建议规范：`novel/chapter/1.json` `xapp/user/xxx.json` `secret/key/xxx.json`

## 🔑 License
`LGPL LollipopKit 2022`