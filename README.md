### 目录结构
- config 所有的配置文件
    - huix.conf.demo 配置文件小样
- main.go 实现的代码文件


### 启动程序
```
$ chmod +x ./backmysql
$ nohup /usr/local/backmysql/./backmysql > /var/log/backmysql.log 2>&1 &
```

### 编译
- `set GOOS=linux/windows`
- `set GOARCH=amd64/arm`
- `go build -o backmysql main.go`


## 注意
修改完配置记得`重新启动程序`!!