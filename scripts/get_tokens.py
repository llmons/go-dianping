import time

import pymysql
import requests

db = pymysql.connect(host="localhost", port=3306, user="root", password="123", database="hmdp")
cursor = db.cursor()

cursor.execute("SELECT phone FROM tb_user LIMIT 1000")
rows = cursor.fetchall()

tokens = []
for (phone,) in rows:
    requests.post(f"http://localhost:8080/api/user/code?phone={phone}")
    time.sleep(0.2)
    resp = requests.post("http://localhost:8080/api/user/login", json={
        "phone": phone,
        "code": "123456"
    })
    try:
        token = resp.json()["data"]["token"]
        tokens.append(token)
    except requests.exceptions.RequestException:
        print(f"Failed for {phone}")
    except ValueError:
        print("Error parsing JSON response:")
    time.sleep(0.2)

with open("tokens.txt", "w") as f:
    for t in tokens:
        f.write(t + "\n")
