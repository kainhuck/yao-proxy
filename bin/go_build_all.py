#!/usr/bin/python3
import os
import sys

os_arch = {
    "darwin": [
        "amd64",
        "arm64"
    ],
    "linux": [
        "386",
        "amd64",
        "arm",
        "arm64",
    ],
    "windows": [
        "386",
        "amd64"
    ]
}
pkg_name = os.getcwd().split("/")[-1]


def build_all(argv: [str], name: str):
    for os_, arch_s in os_arch.items():
        for arch in arch_s:
            if os_ == "windows":
                print(f"CGO_ENABLED=0 GOOS={os_} GOARCH={arch} go build -o {name}_{os_}_{arch}.exe {' '.join(argv)}")
                os.system(f"CGO_ENABLED=0 GOOS={os_} GOARCH={arch} go build -o {name}_{os_}_{arch}.exe {' '.join(argv)}")
            else:
                print(f"CGO_ENABLED=0 GOOS={os_} GOARCH={arch} go build -o {name}_{os_}_{arch} {' '.join(argv)}")
                os.system(f"CGO_ENABLED=0 GOOS={os_} GOARCH={arch} go build -o {name}_{os_}_{arch} {' '.join(argv)}")


def main():
    if len(sys.argv) == 1:
        build_all(sys.argv[1:], pkg_name)
    else:
        if sys.argv[1] != '-o':
            build_all(sys.argv[1:], pkg_name)
        else:
            if len(sys.argv) <= 2:
                print("need output file")
            else:
                build_all(sys.argv[3:], sys.argv[2])


if __name__ == '__main__':
    main()