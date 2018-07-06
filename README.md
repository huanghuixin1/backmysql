### 目录结构
- process-config 所有需要守护的进程配置文件
    - huix.conf.demo 配置文件小样
- main.go 实现的代码文件
- config.conf 程序的配置文件


### 启动程序
```
$ chmod +x ./backmysql
$ nohup ./backmysql > /var/log/backmysql.log 2>&1 &
```