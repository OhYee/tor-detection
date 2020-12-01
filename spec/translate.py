import sys
from imp import reload
import time
import hashlib
import requests
import uuid
import re
import threading
import os

APP_KEY = ""
APP_SECRET = ""


with open(".secret") as f:
    APP_KEY, APP_SECRET = f.read().split("\n")


class Thread(threading.Thread):
    def __init__(self, thread_id, task_handle, task_pool):
        threading.Thread.__init__(self)
        self.thread_id = thread_id
        self.task_handle = task_handle
        self.task_pool = task_pool
        print("Thread", self.thread_id, "start")

    def run(self):
        while 1:
            task = self.task_pool()
            if task == None:
                break
            self.task_handle(self.thread_id, task)

    def __del__(self):
        print("Thread", self.thread_id, "close")


def reformat(text: str) -> str:
    lines = text.split("\n")
    res = ""
    for line in lines:
        if re.match(r'^[ ]*[A-Z0-9\-\*].*$', line):
            res += "\n" + line
        elif re.match(r"^[ ]*$", line):
            res += "\n"
        else:
            res += " " + line.strip(" ")
    return res


reload(sys)


def connect(q: str):
    YOUDAO_URL = 'https://openapi.youdao.com/api'

    def encrypt(signStr):
        hash_algorithm = hashlib.sha256()
        hash_algorithm.update(signStr.encode('utf-8'))
        return hash_algorithm.hexdigest()

    def truncate(q):
        if q is None:
            return None
        size = len(q)
        return q if size <= 20 else q[0:10] + str(size) + q[size - 10:size]

    def do_request(data):
        headers = {'Content-Type': 'application/x-www-form-urlencoded'}
        retry = 10
        response = None
        while retry > 0:
            try:
                response = requests.post(
                    YOUDAO_URL, data=data, headers=headers, timeout=10
                )
                break
            except Exception as e:
                print("Get page error: {}, {} times left...".format(e, retry))
                retry -= 1
        return response

    data = {}
    data['from'] = 'en'
    data['to'] = 'zh-CHS'
    data['signType'] = 'v3'
    curtime = str(int(time.time()))
    data['curtime'] = curtime
    salt = str(uuid.uuid1())
    signStr = APP_KEY + truncate(q) + salt + curtime + APP_SECRET
    sign = encrypt(signStr)
    data['appKey'] = APP_KEY
    data['q'] = q
    data['salt'] = salt
    data['sign'] = sign
    # data['vocabId'] = "您的用户词表ID"

    try:
        response = do_request(data)
        contentType = response.headers['Content-Type']
        result = response.json()
        if result["errorCode"] == "0":
            return result["translation"][0]
        else:
            return q
    except Exception as e:
        print(e)
        return q


def dofile(name: str):
    print(name)
    text = ""
    with open(os.path.join("./torspec", name)) as f:
        text = f.read()

    text = reformat(text)
    with open(os.path.join("./output", name.split(".")[0]+"_format.txt"), "w") as f:
        f.write(text)

    result = text.split("\n")
    l = len(result)
    idx = 0
    finished_task = 0
    task_lock = threading.Lock()
    result_lock = threading.Lock()

    def task_pool():
        nonlocal idx, l
        task_lock.acquire()
        if idx < l:
            res = (idx, result[idx])
            idx += 1
        else:
            res = None
        task_lock.release()
        return res

    def translate(thread_id: int, args: tuple):
        nonlocal finished_task
        idx, text = args
        prefix, text = re.findall(r'^([ ]*)(.*)$', text)[0]
        try:
            temp = prefix + connect(text)
        except:
            temp = prefix + text
        result_lock.acquire()
        result[idx] = temp
        finished_task += 1
        print(finished_task, "/", l, "-", thread_id, idx)
        result_lock.release()

    threads = [
        Thread(i, translate, task_pool)
        for i in range(1024)
    ]
    for thread in threads:
        thread.start()
    for thread in threads:
        thread.join()
    del threads

    with open(os.path.join("./output", name.split(".")[0]+"_zh.txt"), "w", encoding="utf-8") as f:
        f.write("\n".join(result))


if __name__ == "__main__":
    if not os.path.exists("output"):
        os.mkdir("output")
    for root, dirs, files in os.walk("./torspec", topdown=False):
        for name in files:
            if (root == "./torspec" and name.split(".")[-1] == "txt"):
                dofile(name)
