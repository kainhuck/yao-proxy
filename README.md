# yao-proxy



## 介绍

这是一个简单代理，核心代码百来行，便可以绕过防火墙实现访问墙外资源，程序分为本地代理和远程代理，本地代理部署在本地，远程代理部署在墙外可访问的服务器上，详见代码

## 声明

本项目仅用作`学习交流`，`提升编程能力`，请**不要**将该项目用于非法用途!

## 使用

### 快速测试

```
git clone https://github.com/kainhuck/yao-proxy.git
```

```
make run-local
```

```
make run-remote
```

## docker 部署

本地代理：

```
docker pull docker.pkg.github.com/kainhuck/yao-proxy/local:latest

docker run --name yao-proxy --net=host --restart=always -v <your config path>:/etc/yao-proxy/config.json -d docker.pkg.github.com/kainhuck/yao-proxy/local:latest
```

远程代理：

```
docker pull docker.pkg.github.com/kainhuck/yao-proxy/remote:latest

docker run --name yao-proxy --net=host --restart=always -v <your config path>:/etc/yao-proxy/config.json -d docker.pkg.github.com/kainhuck/yao-proxy/remote:latest
```

## 配置文件示例

[local-config](cmd/local/res/config.json)

[remote-config](cmd/remote/res/config.json)

## 贡献代码

`main`分支为最新稳定分支

`develop`分支为最新分支

`release`分支为历史稳定分支，应该从`main`分支切过去

`feature`分支为新特性分支，应该从`develop`中切过去

`fix`分支为bug修复分支



## todo

1. 使用systemd来部署服务

2. 实现cli来安装部署remote，以及生成local的配置文件



## 更新说明

### v2.2.1

- 本地代理更新，可以支持代理多个端口

- 配置文件和之前版本不兼容
