#!/usr/bin/env python3
import os

APP_NAME = 'nano-db'
SYSTEMD_NAME = 'nanodb'


def pull():
    print('[pulling...]')
    code = os.system('git fetch --all')
    if code != 0:
        return False
    code = os.system('git reset --hard')
    if code != 0:
        return False
    code = os.system('git pull')
    return code == 0


def kill():
    print('getting pid...')
    result = os.popen(f"top -cbn1  | grep '{APP_NAME}'")
    if result:
        lines = result.readlines()
        if len(lines) != 2:
            # Skip, because there is no process
            print(f'{APP_NAME} is not running')
            return True
        first_line = lines[0]
        splited = first_line.split(' ')
        pid = splited[0]
        if not pid:
            pid = splited[1]

        print(f'killing pid {pid} ...')
        exitcode = os.system('kill ' + pid)
        return exitcode == 0


def build():
    print('[building...]')
    code = os.system('go build')
    return code == 0


def start():
    print('[restarting...]')
    code = os.system(f'systemctl --user start {SYSTEMD_NAME}')
    return code == 0


if __name__ == '__main__':
    if pull():
        if build():
            if kill():
                if start():
                    print(f'upgrade {APP_NAME} success')
