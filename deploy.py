#!/usr/bin/env python3
import os

APP_NAME = 'nano-db'
SYSTEMD_NAME = 'nanodb'


def pull():
    print('[pulling...]')
    code = os.system('git fetch --all')
    if code != 0:
        return False
    code = os.system('git reset --hard origin/main')
    if code != 0:
        return False
    code = os.system('git pull')
    return code == 0


def build():
    print('[building...]')
    code = os.system('/usr/local/go/bin/go build')
    return code == 0


def restart():
    print('[restarting...]')
    code = os.system(f'systemctl --user start {SYSTEMD_NAME}')
    return code == 0


if __name__ == '__main__':
    if pull():
        if build():
            if restart():
                print(f'upgrade {APP_NAME} success')
