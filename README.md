# 简介
这是一个使用Golang实现的代理服务器程序，支持HTTP和SOCKS5代理协议。该代理服务器可以在本地搭建一个HTTP和SOCKS5代理服务器，方便用户在需要的时候进行代理访问。

# 功能特性
- 支持HTTP和SOCKS5代理协议；
- 支持匿名代理和带有账号密码认证的代理；
- 支持SOCKS5代理的TCP协议的转发；
# 运行
1. 安装Go编程语言，可以参考官网的[安装指南](https://golang.org/doc/install)进行安装；
2. 下载并安装本项目源码；
3. 在命令行中进入项目目录，运行以下命令：
```
go mod tidy
```
4. 启动代理服务器 
```
go run cmd/main.go
```