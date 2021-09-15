#!/usr/bin/python3
import os
import sys


def push(msg: str):
    print("git add .")
    os.system("git add .")
    print(f'git commit -m "{msg}"')
    os.system(f'git commit -m "{msg}"')
    print("git push")
    os.system("git push")


if __name__ == '__main__':
    if len(sys.argv) == 1:
        push("add .")
    else:
        push(" ".join(sys.argv[1:]))