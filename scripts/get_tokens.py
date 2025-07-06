import time
import requests
import pymysql

db = pymysql.connect(host="localhost",port=3306, user="root", password="123", database="hmdp")
cursor = db.cursor()

cursor.execute("SELECT phone FROM tb_user LIMIT 1000")
rows = cursor.fetchall()

tokens = []
for (phone,) in rows:
    requests.post(f"http://localhost:8080/api/user/code?phone={phone}")
    resp = requests.post("http://localhost:8080/api/user/login", json={
        "phone": phone,
        "code": "123456"
    })
    try:
        token = resp.json()["data"]["token"]
        tokens.append(token)
    except:
        print(f"Failed for {phone}")
    time.sleep(0.1)

with open("tokens.txt", "w") as f:
    for t in tokens:
        f.write(t + "\n")
