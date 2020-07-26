- [x] golang
- [x] buffalo
- [x] send email
- [x] receive email

#### 仅支持 `QQ`邮箱的 `smtp`协议和 `imap`协议
#### 仅支持 加密方式登录,使用授权码登录

#### 支持使用默认账号免登录发送信息
#### 退出后删除所有信息,不保留用户隐私信息
##### 使用方法,clone后新建.env文件,配置数据库及默认邮箱的账号密码

### 使用方法
1. clone 下来,启动需要 `buffalo`的`cli`程序,[buffalo](https://github.com/gobuffalo/buffalo)
2. 配置数据库,将下面配置写入`env`文件中,
```
    DATABASE=xxx
    USER=xxx
    PASSWORD=xxx
    HOST=xxx
```
3. 创建数据库 ` buffalo db create -a -d` 
4. 迁移数据库 `buffalo db migrate up` 
5. 配置 邮箱服务器,也是在 `env`  文件中
```
SMTP_PORT=587
SMTP_HOST=smtp.qq.com
SMTP_USER=xxx
SMTP_PASSWORD=xxx

IMAP_HOST=imap.qq.com
IMAP_PORT=993
```
提前生成了一个`.env.example`文件,更改文件名为`.env`

配置完成后使用 `buffalo dev` 即可启动服务,默认是 3000端口,可在.`.env`文件中更改端口号

