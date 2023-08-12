## Nano DB
一款以golang编写的轻量、非关系型、基于文件系统的数据库。  

白话文：将数据按文件储存，再提供http接口来访问。  

## 🔖 特点
- 轻量：即使包含数十万索引，树莓派上也能流畅运行
- 高速：微秒级查询
- RESTful接口：HTTP协议，方便使用
- 权限管理：ACL
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
[INF] generated token: FHYmGdNwfiJngvF2z
[SUC] acl update rule: success
```
可以用生成的 `token` 访问名为 `novel` 的数据库

⚠️ **注意**：
- 请在使用http接口时，在 `headers` 内的 `NanoDB` 键附带此 `token`
- 使用 sdk，需要用到此 `token` 

#### 添加权限

##### 手动添加
建议手动添加，这样不会将 `token` 暴露至 shell，可以打开 `.cfg/acl.json`（如文件不存在，需要先启动数据库一次或手动创建）文件进行手动修改，例如：
```json
{
      "ver": 1,
      "rules":[
            {
                  "token": "token1",
                  "dbs":["novel"]
            }
      ]
}
```

如果想给 `token1` 添加 `test` 数据库的权限，可以如下修改：
```json
{
      "ver": 1,
      "rules":[
            {
                  "token": "token1",
                  "dbs":["novel", "test"]
            }
      ]
}
```

##### 命令行添加
`nano-db -d {dbName} -t {token}`   

示例：
```
➜  nano-db git:(main) ✗ nano-db -d novel -t FHYmGdNwfiJngvF2z
[SUC] acl update rule: success
```
然后就可以用该 `token` 访问名叫 `novel` 的数据库了

⚠️**注意**：
如果当前数据库正在运行，acl更改将在一分钟内应用。

### 🔨 数据库操作
接下来是 http 接口的使用。

方法|接口|功能|额外说明
---|---|---|---
HEAD|`/`|查看数据库是否存活|唯一不需要 token 的接口，可用于客户端检查数据库是否存活
GET|`/{DB}`|获取 DB 内所有 DIR|会返回所有 DIR 的**名称**
DELETE|`/{DB}`|删除数据库|不会删除对该数据库的权限
GET|`/{DB}/{DIR}`|获取 DIR 内所有 FILE|获取 DIR 下所有文件的**名称**
DELETE|`/{DB}/{DIR}`|删除某 DIR|并删除 DIR 下所有 FILE
GET|`/{DB}/{DIR}/{FILE}`|获取 FILE|不存在则会返回错误
POST|`/{DB}/{DIR}/{FILE}`|更新 FILE|需要在 body 附带需要写入的数据
DELETE|`/{DB}/{DIR}/{FILE}`|删除|如果路径不存在则会返回错误

⚠️**注意**：`{DB}`,`{DIR}`,`{FILE}` 不能包含除 `.`(不能在开头结尾) `0-9` `A-Z` `a-z` `_` `-` 以外的字符，并且他们的长度都不能超过37（可在配置中自定义）。

建议规范：`novel/chapter/1.json` `xapp/user/xxx.json` `secret/key/xxx.json`

## 🔑 License
`LGPL LollipopKit 2022`