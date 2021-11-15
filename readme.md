### 框架
```
beego [https://github.com/astaxie/beego]
```

### 中间件
```

关系数据:mysql:
beego自带orm
安装:
sudo docker run -d --restart=always -p 0.0.0.0:3306:3306 --name root -e MYSQL_ROOT_PASSWORD=123456 -d mysql
文件存储:minio:
https://github.com/minio/minio
安装:
docker run -p 0.0.0.0:9000:9000 -p 0.0.0.0:9001:9001 --name minio \
-d --restart=always \
-e "MINIO_ACCESS_KEY=minio" \
-e "MINIO_SECRET_KEY=putao520" \
-v /home/data:/data \
-v /home/config:/root/.minio \
minio/minio server /data --console-address "0.0.0.0:9001"
http://192.168.75.80:9001
```