# yao-proxy



## 介绍

这是一个简单代理，核心代码百来行，便可以绕过防火墙实现访问墙外资源，程序分为本地代理和远程代理，本地代理部署在本地，远程代理部署在墙外可访问的服务器上，详见代码

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
docker run --name yp-proxy --net=host -v <your config path>:/etc/yao-proxy/config.json -d github.com/kainhuck/yao-proxy-local:2.1.1 
```

远程代理：

```
docker run --name yp-proxy --net=host -v <your config path>:/etc/yao-proxy/config.json -d github.com/kainhuck/yao-proxy-remote:2.1.1 
```

## 贡献代码

`main`分支为最新稳定分支

`develop`分支为最新分支

`release`分支为历史稳定分支，应该从`main`分支切过去

`feature`分支为新特性分支，应该从`develop`中切过去

`fix`分支为bug修复分支



## todo

增加docker部署方式 👌🏻

