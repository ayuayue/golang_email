- [x] golang
- [x] buffalo
- [x] send email
- [x] receive email

#### 仅支持 `QQ`邮箱的 `smtp`协议和 `imap`协议

#### 仅支持 加密方式登录,使用授权码登录

#### 支持使用默认账号免登录发送信息

#### 不使用数据库,仅使用 `session` 功能和缓存进行邮件收发

#### 退出后删除所有信息,不保留用户隐私信息

### 使用方法 1 源码运行

1. clone 下来,启动需要 `buffalo`的`cli`程序,[buffalo](https://github.com/gobuffalo/buffalo) 
2. 如想开启免登录,请先到 `.env` 文件中配置系统默认发件邮箱账号和授权码
3. 配置完成后使用 `buffalo dev` 即可启动服务,默认是 3000 端口,可在.`.env`文件中更改端口号
### 二进制文件运行

1. 下载 release 页中的二进制文件
2. 会自动生成`.env` 文件,和上面一样进行配置,
3. 授予执行权限,终端运行

#### 路由
运行以后,在跟路由下的 `api` 路径可以查看所有路由 如 `localhost:3000/api`