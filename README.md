# yao-proxy

![GitHub](https://img.shields.io/github/license/kainhuck/yao-proxy) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kainhuck/yao-proxy) ![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/kainhuck/yao-proxy)

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

**注意📢: docker镜像不再发布到github packages，现在只发布到dockerhub：[🔗](https://hub.docker.com/repository/docker/kainhuck/yao-proxy)**

现在将两个镜像发布到同一个仓库，通过tag来区分，

- local镜像tag

  latest-local

- remote镜像tag

  latest-remote

部署方式如下：

_注意：_

_1. 运行时请指定配置文件的路径，[配置文件示例](#配置文件示例)_ 

_2. mac系统不支持host模式，请手动通过 -p 来映射端口_

**本地代理：**

```shell
docker run --name yao-proxy \
           --net=host 
           --restart=always 
           -v <your config path>:/etc/yao-proxy/config.json \
           -d kainhuck/yao-proxy:latest-local
```

**远程代理：**

```shell
docker run --name yao-proxy \
           --net=host \
           --restart=always \
           -v <your config path>:/etc/yao-proxy/config.json \
           -d kainhuck/yao-proxy:latest-remote
```

## 二进制部署

1. 下载最新的对应平台的二进制文件：[🔗](https://github.com/kainhuck/yao-proxy/releases)

2. 准备好配置文件

3. 运行程序 `-c` 指定配置文件，例:

   ```
   ./local_darwin_amd64 -c /etc/yao-proxy/config.json
   ./remote_darwin_amd64 -c /etc/yao-proxy/config.json
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

### v2.2.3

- 过滤规则增加ipv4区间写法，参考[local-config](cmd/local/res/config.json)

### v2.2.2

- 本地代理新增过滤规则，可以不代理指定的域名或者IP地址，写法参考[local-config](cmd/local/res/config.json)

### v2.2.1

- 本地代理更新，可以支持代理多个端口

- 配置文件和之前版本不兼容
