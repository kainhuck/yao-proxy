#!/usr/bin/python3
import os
import sys

BUILD_CMD_FORMAT = 'docker buildx build --platform {platform} -t {url}:{version}-{name} -f cmd/{name}/Dockerfile ' \
                   '--push . '


def check_ok(name: str, code: int):
    if code != 0:
        print(f" ############   build `{name}` fail!   ############\n")
    else:
        print(f" ############   build `{name}` success!   ############\n")


def build(platform: str, url: str, version: str, name: str):
    build_cmd = BUILD_CMD_FORMAT.format(platform=platform, url=url, version=version, name=name)
    print(build_cmd)
    check_ok(name, os.system(build_cmd))


def main():
    url = "kainhuck/yao-proxy"
    name = ""
    platform = ""
    version = ""
    if len(sys.argv) == 1:
        print("缺少运行参数：./bin/docker-build.py <*服务名(all/local/remote)> <平台(linux/amd64,linux/arm/v7,linux/arm64)> <版本("
              "1.0.0)>")
        return

    name = sys.argv[1]

    if len(sys.argv) > 2:
        platform = sys.argv[2]

    if len(sys.argv) > 3:
        version = sys.argv[3]

    if platform == "":
        print("未指定平台默认打包: linux/amd64,linux/arm/v7,linux/arm64")
        platform = "linux/amd64,linux/arm/v7,linux/arm64"

    if version == "":
        print("未指定版本号默认为: latest")
        version = "latest"

    if name == "all":
        build(platform, url, version, "local")
        build(platform, url, version, "remote")
    elif name == "local":
        build(platform, url, version, "local")
    elif name == "remote":
        build(platform, url, version, "remote")
    else:
        print("不支持的服务: " + name)


if __name__ == '__main__':
    main()
